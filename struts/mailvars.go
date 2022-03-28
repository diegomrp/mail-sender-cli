package struts

type MailVars struct {
	Host         string
	Port         int
	Template     string
	PathImage    string
	PathDocument string
	From         string
	Subject      string
	RecipientsFile string 
}
