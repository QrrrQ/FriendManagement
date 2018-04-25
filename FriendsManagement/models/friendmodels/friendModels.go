package friendmodels

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql" //init mysql driver before  package init
)

type User struct {
	UserId   uint64
	Email    string
	Password string
	Deleted  bool //0 normal  1 delete
}

type Connection struct {
	Index    uint64
	User1    uint64 //requestor
	User2    uint64 //target
	Relation int    //0 无  1 朋友 2 申请中（user1为申请人，user2为被申请人）3 block（user1为block发起人） 4 删除
	Deleted  bool
}

var db = &sql.DB{}

func init() {
	db, _ = sql.Open("mysql", "root:123456@/qrrrq?charset=utf8")
	// if err != nil {
	// 	panic("can not open database.")
	// }
	// fmt.Println("db open success")
}

func (c *User) CreateUser() error {
	_, err := db.Exec(`INSERT INTO user (email,password) VALUES(?,?)`, c.Email, c.Password)
	return err
}

func (c *User) CheckUserEmailExist() bool {
	var email string
	db.QueryRow("SELECT email FROM user WHERE email = ? AND deleted = 0", c.Email).Scan(&email)
	if email == "" {
		return false
	}
	return true
}

func (c *User) GetUserByEmail() error {
	err := db.QueryRow("SELECT * FROM user WHERE email = ? AND deleted = 0", c.Email).Scan(&c.UserId, &c.Email, &c.Password, &c.Deleted)
	fmt.Println(err)
	fmt.Println("GetUserByEmail")
	if err != nil {
		return err
	}
	if c.Email == "" {
		err = errors.New("no user info")
	}
	return err
}

func (c *Connection) CreateOrUpdateConnection() error {
	_, err := db.Exec(`INSERT INTO connection (user1,user2,relation,deleted) 
	values(?,?,?,?) ON DUPLICATE KEY UPDATE relation=VALUES(relation), deleted=VALUES(deleted)`, c.User1, c.User2, c.Relation, 0)
	return err
}

func (c *Connection) GetRelationByUserIds(uid1, uid2 uint64) error {
	err := db.QueryRow("SELECT * FROM connection WHERE user1 = ? AND user2 = ? AND deleted = 0", uid1, uid2).Scan(&c.Index, &c.User1, &c.User2, &c.Relation, &c.Deleted)
	return err
}

func (c *User) GetFriendEmailList() []string {
	var v []string
	rows, err := db.Query(`SELECT a.email FROM user a WHERE a.user_id IN (SELECT b.user2 FROM connection b 
	WHERE b.user1 = ? and b.relation = 1 and b.deleted = 0 order by a.email)`, c.UserId)
	if err != nil {
		fmt.Println("GetFriendEmailList get error 1")
	}

	for rows.Next() {
		var str string
		if err = rows.Scan(&str); err != nil {
			fmt.Println("GetFriendEmailList scan error")
		}
		v = append(v, str)
	}

	return v
}

func (c *User) GetCommonList(user2 User) []string {
	var v []string
	rows, _ := db.Query(`SELECT a.email FROM user a WHERE a.deleted = 0 AND a.user_id IN (SELECT b.user2 FROM connection b 
		WHERE b.user1 = ? AND b.deleted = 0 AND b.user2 IN (SELECT c.user2 FROM connection c 
			 WHERE c.user1 = ? AND deleted = 0))`, c.UserId, user2.UserId)
	for rows.Next() {
		var str string
		if err := rows.Scan(&str); err != nil {
			fmt.Println("GetCommonList scan error")
		}
		v = append(v, str)
	}
	return v
}

func (c *Connection) CheckRalation() error {
	err := db.QueryRow("SELECT * FROM connection WHERE user1 = ? AND user2 = ? AND deleted = 0", c.User1, c.User2).Scan(&c.Index, &c.User1, &c.User2, &c.Relation, &c.Deleted)
	return err
}

func (c *User) GetAvailabelEmails() []string {
	var v []string
	rows, _ := db.Query(`
		SELECT a.email FROM user a WHERE a.user_id =(SELECT b.user1 FROM connection b WHERE 
			b.user2 = ? AND b.deleted = 0 AND (b.relation = 1 OR b.relation = 2) AND b.user1 NOT IN 
			 (SELECT c.user2  FROM connection c WHERE 
		c.user1 = ? AND c.deleted = 0 AND (c.relation = 3 OR c.relation = 4)))
		`, c.UserId, c.UserId)
	for rows.Next() {
		var str string
		if err := rows.Scan(&str); err != nil {
			fmt.Println("GetCommonList scan error")
		}
		v = append(v, str)
	}
	return v
}
