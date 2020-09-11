package models

type ChangeBalanceReq struct {
	UserId    int       `json:"user_id"`
	Change    int       `json:"change"`
}

type ChangeBalanceResp struct {
	UserId    int       `json:"user_id"`
	Change    int       `json:"change"`
	Message   string	`json:"message"`
}

type TransferReq struct {
	UserId    int       `json:"user_id"`
	Change    int       `json:"change"`
	TargetId  int       `json:"target_id"`
}

type TransferResp struct {
	UserId    int       `json:"user_id"`
	Change    int       `json:"change"`
	TargetId  int       `json:"target_id"`
	Message   string	`json:"message"`
}

type GetBalanceReq struct {
	UserId    int       `json:"user_id"`
}

type GetBalanceResp struct {
	UserId    int       `json:"user_id"`
	Balance    int      `json:"balance"`
}

type GetTransactionsReq struct {
	UserId    			int         `json:"user_id"`
	TransactionsOnPage 	int		    `json:"transactions_on_page"`
	ChangeSort			bool		`json:"change_sort"`
	TimeSort			bool		`json:"time_sort"`
}

type Transaction struct {
	Change   	int					`json:"change"`
	ChangeTime 	string				`json:"change_time"`
	Comment     string				`json:"comment"`
}

type GetTransactionsResp struct {
	UserId    			int         	`json:"user_id"`
	CurrentBalance    	int         	`json:"current_balance"`
	Transactions		[][]Transaction `json:"transactions"`
}



