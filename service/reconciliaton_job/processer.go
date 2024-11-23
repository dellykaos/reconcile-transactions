package reconciliatonjob

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/delly/amartha/entity"
	filestorage "github.com/delly/amartha/repository/file_storage"
	dbgen "github.com/delly/amartha/repository/postgresql"
)

type bankMappedTransactions struct {
	bankName     string
	transactions map[string][]*entity.Transaction
}

// Processer is a contract to process pending reconciliation job
type Processer interface {
	Process(ctx context.Context) error
}

// ProcesserRepository is a dependency of repository that needed to process reconciliation job
type ProcesserRepository interface {
	ListPendingReconciliationJobs(ctx context.Context) ([]dbgen.ReconciliationJob, error)
	SaveFailedReconciliationJob(ctx context.Context, arg dbgen.SaveFailedReconciliationJobParams) (dbgen.ReconciliationJob, error)
	SaveSuccessReconciliationJob(ctx context.Context, arg dbgen.SaveSuccessReconciliationJobParams) (dbgen.ReconciliationJob, error)
}

// FileGetter is a dependency of repository that needed to get file from storage
type FileGetter interface {
	Get(ctx context.Context, filePath string) (*filestorage.File, error)
}

// ProcesserService is an implementation of Processer to process
// pending reconciliation job
type ProcesserService struct {
	repo    ProcesserRepository
	storage FileGetter
}

var _ = Processer(&ProcesserService{})

// NewProcesserService create new processer service
func NewProcesserService(repo ProcesserRepository, storage FileGetter) *ProcesserService {
	return &ProcesserService{repo: repo, storage: storage}
}

// Process process pending reconciliation job
func (s *ProcesserService) Process(ctx context.Context) error {
	jobs, err := s.getPendingReconciliationJobs(ctx)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		if err = s.processReconciliationJob(ctx, job); err != nil {
			job.Status = entity.ReconciliationJobStatusFailed
			job.ErrorInformation = err.Error()
		}
		if job.Status == entity.ReconciliationJobStatusFailed {
			if s.saveFailedJob(ctx, job); err != nil {
				fmt.Printf("failed to update job: %d status to failed with error: %v\n", job.ID, err)
			}
		} else if job.Status == entity.ReconciliationJobStatusSuccess {
			if err = s.saveSuccessJob(ctx, job); err != nil {
				fmt.Printf("failed to save success job: %d with error: %v\n", job.ID, err)
			}
		}
	}

	return nil
}

func (s *ProcesserService) saveSuccessJob(ctx context.Context, job *entity.ReconciliationJob) error {
	params := dbgen.SaveSuccessReconciliationJobParams{
		ID: job.ID,
	}
	params.Result.Set(job.Result)
	if _, err := s.repo.SaveSuccessReconciliationJob(ctx, params); err != nil {
		return err
	}

	return nil
}

func (s *ProcesserService) saveFailedJob(ctx context.Context, job *entity.ReconciliationJob) error {
	if _, err := s.repo.SaveFailedReconciliationJob(ctx, dbgen.SaveFailedReconciliationJobParams{
		ID:               job.ID,
		ErrorInformation: sql.NullString{String: job.ErrorInformation, Valid: true},
	}); err != nil {
		return err
	}

	return nil
}

func (s *ProcesserService) processReconciliationJob(ctx context.Context, job *entity.ReconciliationJob) error {
	systemTrxFile, bankFiles, err := s.getCSVFiles(ctx, job)
	if err != nil {
		return err
	}

	systemTrx, err := s.convertSystemCSVFileToTransactions(job, systemTrxFile)
	if err != nil {
		return err
	}

	bankTrxs := []*bankMappedTransactions{}
	for bankName, file := range bankFiles {
		trx, err := s.convertBankCSVFileToMappedTransactions(job, file)
		if err != nil {
			return err
		}
		bankTrxs = append(bankTrxs, &bankMappedTransactions{
			bankName:     bankName,
			transactions: trx,
		})
	}

	result := s.processReconciliation(job, systemTrx, bankTrxs)
	job.Result = result
	job.Status = entity.ReconciliationJobStatusSuccess

	return nil
}

func (s *ProcesserService) processReconciliation(job *entity.ReconciliationJob,
	systemTrx []*entity.Transaction, bankTrxs []*bankMappedTransactions) *entity.ReconciliationResult {
	result := &entity.ReconciliationResult{
		TotalTransactionProcessed: 0,
		TotalTransactionMatched:   0,
		TotalTransactionUnmatched: 0,
		TotalDiscrepancyAmount:    0,
		MissingTransactions:       []entity.Transaction{},
		MissingBankTransactions:   map[string][]entity.Transaction{},
	}

	for _, trx := range systemTrx {
		var found bool
		var bankTrxIdx int
		date := trx.Time.Format(time.DateOnly)
		result.TotalTransactionProcessed++
		for _, bankTrx := range bankTrxs {
			if bankTrxTrxs, ok := bankTrx.transactions[date]; ok {
				for idx, bankTrx := range bankTrxTrxs {
					discrepancyThreshold := float64(job.DiscrepancyThreshold) * trx.Amount
					minDiscrepancy := trx.Amount - discrepancyThreshold
					maxDiscrepancy := trx.Amount + discrepancyThreshold
					bankAmountInThreshold := bankTrx.Amount >= minDiscrepancy && bankTrx.Amount <= maxDiscrepancy
					if bankAmountInThreshold {
						bankTrxIdx = idx
						found = true
						break
					}
				}
				if found {
					result.TotalTransactionProcessed++
					result.TotalTransactionMatched++
					// remove matched bank transaction, so it won't be processed again
					// and we can track missing bank transactions
					bankTrx.transactions[date] = append(bankTrx.transactions[date][:bankTrxIdx],
						bankTrx.transactions[date][bankTrxIdx+1:]...)
					break
				}
			}
		}
		if !found {
			// add missing system transaction to result
			result.MissingTransactions = append(result.MissingTransactions, *trx)
			result.TotalDiscrepancyAmount += trx.Amount
			result.TotalTransactionUnmatched++
			continue
		}
	}

	for _, bankTrx := range bankTrxs {
		for _, trxs := range bankTrx.transactions {
			for _, trx := range trxs {
				result.MissingBankTransactions[bankTrx.bankName] = append(result.MissingBankTransactions[bankTrx.bankName], *trx)
				result.TotalDiscrepancyAmount += trx.Amount
				result.TotalTransactionUnmatched++
				result.TotalTransactionProcessed++
			}
		}
	}

	return result
}

func (s *ProcesserService) convertSystemCSVFileToTransactions(job *entity.ReconciliationJob,
	file *filestorage.File) ([]*entity.Transaction, error) {
	csvReader := csv.NewReader(file.Buf)
	result := []*entity.Transaction{}
	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		trx, err := s.convertSystemTransactionRecordToTransaction(record)
		if err != nil {
			return nil, err
		}

		notInRange := trx.Time.Before(job.StartDate) || trx.Time.After(job.EndDate)
		if notInRange {
			continue
		}

		result = append(result, trx)
	}

	return result, nil
}

func (s *ProcesserService) convertBankCSVFileToMappedTransactions(job *entity.ReconciliationJob,
	file *filestorage.File) (map[string][]*entity.Transaction, error) {
	csvReader := csv.NewReader(file.Buf)
	result := map[string][]*entity.Transaction{}
	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		trx, err := s.convertBankTransactionRecordToTransaction(record)
		if err != nil {
			return nil, err
		}

		notInRange := trx.Time.Before(job.StartDate) || trx.Time.After(job.EndDate)
		if notInRange {
			continue
		}

		date := trx.Time.Format(time.DateOnly)
		if _, ok := result[date]; !ok {
			result[date] = []*entity.Transaction{}
		}
		result[date] = append(result[date], trx)
	}

	return result, nil
}

func (s *ProcesserService) convertSystemTransactionRecordToTransaction(record []string) (*entity.Transaction, error) {
	trxID := record[0]
	amount, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return nil, err
	}
	trxType := entity.TransactionType(record[2])
	transactionTime, err := time.Parse(time.RFC3339, record[3])
	if err != nil {
		return nil, err
	}

	return &entity.Transaction{
		ID:     trxID,
		Amount: amount,
		Type:   trxType,
		Time:   transactionTime,
	}, nil
}

func (s *ProcesserService) convertBankTransactionRecordToTransaction(record []string) (*entity.Transaction, error) {
	trxID := record[0]
	amount, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return nil, err
	}
	var trxType entity.TransactionType
	if amount < 0 {
		amount = amount * -1
		trxType = entity.TxTypeDebit
	} else {
		trxType = entity.TxTypeCredit
	}
	transactionTime, err := time.Parse(time.DateOnly, record[2])
	if err != nil {
		return nil, err
	}

	return &entity.Transaction{
		ID:     trxID,
		Amount: amount,
		Type:   trxType,
		Time:   transactionTime,
	}, nil
}

func (s *ProcesserService) getCSVFiles(ctx context.Context, job *entity.ReconciliationJob) (systemTrxFile *filestorage.File, bankFiles map[string]*filestorage.File, err error) {
	systemTrxFile, err = s.storage.Get(ctx, job.SystemTransactionCsvPath)
	if err != nil {
		return nil, nil, err
	}

	bankFiles = map[string]*filestorage.File{}
	for _, bankFile := range job.BankTransactionCsvPaths {
		file, err := s.storage.Get(ctx, bankFile.FilePath)
		if err != nil {
			return nil, nil, err
		}
		bankFiles[bankFile.BankName] = file
	}

	return systemTrxFile, bankFiles, nil
}

func (s *ProcesserService) getPendingReconciliationJobs(ctx context.Context) ([]*entity.ReconciliationJob, error) {
	jobs, err := s.repo.ListPendingReconciliationJobs(ctx)
	if err != nil {
		return nil, err
	}

	result := []*entity.ReconciliationJob{}
	for _, job := range jobs {
		result = append(result, convertToEntityReconciliationJob(job))
	}

	return result, nil
}
