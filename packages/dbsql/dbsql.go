package dbsql

/* sql statements */
var (
	UserLogin                                                                                       string
	AddUser, GetAllUsersByStatus, UpdateLastlogin, UpdateUser, CountUsersByStatus                   string
	AddMessage, GetAllMessagesByStatus, GetMessageContent, UpdateMessageStatus, CountUnreadMessages string
)

func init() {
	// Logins
	UserLogin = `SELECT name, password, title, lastlogin FROM users WHERE email=? AND status=1`
	// Users
	AddUser = `INSERT INTO users (name,password,company,email,phone,url,picture) VALUES (?, ?, ?, ?, ?, ?, ?)`
	GetAllUsersByStatus = `SELECT name, title, password, company, email, phone, url, comment, picture, lastlogin, status, qturhm, created FROM users WHERE status LIKE ? ORDER BY status ASC, lastlogin ASC`
	UpdateLastlogin = `UPDATE users SET lastlogin=NOW() WHERE email=?`
	UpdateUser = `UPDATE users SET %s=? WHERE email=?`
	CountUsersByStatus = `SELECT COUNT(email) FROM users WHERE status LIKE ?`
	// Messages
	AddMessage = `INSERT INTO messages (user,name,company,email,phone,url,message) VALUES (?, ?, ?, ?, ?, ?, ?)`
	GetAllMessagesByStatus = `SELECT id, user, name, company, email, phone, url, message, status, qturhm, created FROM messages WHERE status LIKE ? ORDER BY status ASC, created ASC`
	GetMessageContent = `SELECT message FROM messages WHERE email=?`
	UpdateMessageStatus = `UPDATE messages SET status=? WHERE id=?`
	CountUnreadMessages = `SELECT COUNT(id) FROM messages WHERE status=0`
}
