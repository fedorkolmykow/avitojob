package models

import "time"

type ChangeBalanceReq struct {
	UserId    int       `json:"user_id"`
	Change    float64   `json:"change"`
	Comment   string	`json:"comment"`
	Source    string    `json:"source"`
}

type ChangeBalanceResp struct {
	UserId    int       `json:"user_id"`
	Balance   float64   `json:"balance"`
}

type TransferReq struct {
	UserId    int       `json:"user_id"`
	Change    float64   `json:"change"`
	TargetId  int       `json:"target_id"`
	Comment   string	`json:"comment"`
}

type TransferResp struct {
	Source ChangeBalanceResp	`json:"source"`
	Target ChangeBalanceResp	`json:"target"`
}

type GetBalanceReq struct {
	UserId    int       `json:"user_id"`
}

type GetBalanceResp struct {
	UserId    int       `json:"user_id"`
	Balance   float64   `json:"balance"`
}

type GetTransactionsReq struct {
	UserId    			int         `json:"user_id"`
	Page				int			`json:"page"`
	TransactionsOnPage 	int		    `json:"transactions_on_page"`
	ChangeSort			bool		`json:"change_sort"`
	TimeSort			bool		`json:"time_sort"`
}

//type User struct{
//	UserId    int       `json:"user_id"`
//	Balance   float64   `json:"balance"`
//}

type Transaction struct {
	UserId          int                 `json:"-" db:"user_id"`
	InitialBalance  float64				`json:"init_balance" db:"init_balance"`
	Change   		float64				`json:"change" db:"change"`
	ChangeTime 		string				`json:"change_time" db:"time"`
	Source          string              `json:"source" db:"source"`
	Comment     	string				`json:"comment" db:"comment"`
}

type Transactions struct{
	Transactions		[]Transaction
	ChangeSort			bool
	TimeSort			bool
}

type GetTransactionsResp struct {
	UserId    			int         	`json:"user_id"`
	CurrentBalance    	float64         `json:"current_balance"`
	Transactions		[]Transaction   `json:"transactions"`
}


func (t Transactions) Less(i, j int) bool {
	if t.ChangeSort{
		return t.Transactions[i].Change < t.Transactions[j].Change
	}
	if t.TimeSort{
		ti, _ := time.Parse(time.RFC822 ,t.Transactions[i].ChangeTime)
		tj, _ := time.Parse(time.RFC822 ,t.Transactions[j].ChangeTime)
		return ti.Before(tj)
	}
	return false
}

func (t Transactions) Swap(i, j int) {
	t.Transactions[i], t.Transactions[j] = t.Transactions[j], t.Transactions[i]
}

func (t Transactions) Len() int {
	return len(t.Transactions)
}
