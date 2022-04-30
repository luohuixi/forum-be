package main

import (
	"forum-user/pkg/auth"
	pb "forum-user/proto"
	s "forum-user/service"
	"forum/config"
	logger "forum/log"
	"forum/model"
	"forum/pkg/handler"
	tracer "forum/pkg/tracer"
	"github.com/micro/go-micro"
	"github.com/opentracing/opentracing-go"
	"log"

	// _ "github.com/micro/go-plugins/registry/kubernetes"

	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/spf13/viper"
)

func main() {
	// init config
	if err := config.Init("", "FORUM_USER"); err != nil {
		panic(err)
	}

	t, io, err := tracer.NewTracer(viper.GetString("local_name"), viper.GetString("tracing.jager"))
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	defer logger.SyncLogger()

	// set var t to Global Tracer (opentracing single instance mode)
	opentracing.SetGlobalTracer(t)

	// init db
	model.DB.Init()
	defer model.DB.Close()

	// init oauth-manager and some variables
	auth.InitVar()
	auth.OauthManager.Init()

	srv := micro.NewService(
		micro.Name(viper.GetString("local_name")),
		micro.WrapHandler(
			opentracingWrapper.NewHandlerWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapHandler(handler.ServerErrorHandlerWrapper()),
	)

	// Init will parse the command line flags.
	srv.Init()

	// Register handler
	pb.RegisterUserServiceHandler(srv.Server(), &s.UserService{})

	// Run the server
	if err := srv.Run(); err != nil {
		logger.Error(err.Error())
	}
}

//
// // FileModel ... 文件物理模型
// type FileModel struct {
// 	ID        uint32 `json:"id" gorm:"column:id;not null" binding:"required"`
// 	URL       string `json:"url" gorm:"column:url;" binding:"required"`
// 	Re        bool   `json:"re" gorm:"column:re;" binding:"required"`
// 	ProjectID uint32 `json:"projectId" gorm:"column:project_id;" binding:"required"`
// 	FatherId  uint32 `json:"father_id" gorm:"column:father_id;" binding:"required"`
// }
//
// // DocModel ... 文档物理模型
// type DocModel struct {
// 	ID        uint32 `json:"id" gorm:"column:id;not null" binding:"required"`
// 	Content   string `json:"content" gorm:"column:content;" binding:"required"`
// 	Re        bool   `json:"re" gorm:"column:re;" binding:"required"`
// 	ProjectID uint32 `json:"projectId" gorm:"column:project_id;" binding:"required"`
// 	FatherId  uint32 `json:"father_id" gorm:"column:father_id;" binding:"required"`
// }
//
// func main() {
//
// }
// func CreateDoc(db *gorm.DB, doc *DocModel) (uint32, error) {
// 	tx := db.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()
//
// 	if err := tx.Create(doc).Error; err != nil {
// 		tx.Rollback()
// 		return 0, err
// 	}
//
// 	isFatherProject := false
// 	fatherId := doc.FatherId
// 	if doc.FatherId == 0 {
// 		isFatherProject = true
// 		fatherId = doc.ProjectID
// 	}
//
// 	if err := AddDocChildren(tx, isFatherProject, fatherId, doc); err != nil {
// 		tx.Rollback()
// 		return 0, err
// 	}
//
// 	return doc.ID, tx.Commit().Error
// }
//
// func CreateFile(db *gorm.DB, file *FileModel) (uint32, error) {
// 	tx := db.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()
//
// 	if err := tx.Create(file).Error; err != nil {
// 		tx.Rollback()
// 		return 0, err
// 	}
//
// 	isFatherProject := false
// 	fatherId := file.FatherId
// 	if file.FatherId == 0 {
// 		isFatherProject = true
// 		fatherId = file.ProjectID
// 	}
//
// 	if err := AddFileChildren(tx, isFatherProject, fatherId, file); err != nil {
// 		tx.Rollback()
// 		return 0, err
// 	}
//
// 	return file.ID, tx.Commit().Error
// }
//
// func addChildren(children string, id uint32) (string, error) {
// 	return "", nil
// }
//
// type Item struct {
// 	Children string
// }
//
// func GetFolder(fatherId uint32) Item {
// 	return Item{}
// }
//
// // AddDocChildren ... 新增 doc 文件树
// func AddDocChildren(tx *gorm.DB, isFatherProject bool, fatherId uint32, doc *DocModel) error {
// 	id := doc.ID
//
// 	item := GetFolder(fatherId)
//
// 	newChildren, err := addChildren(item.Children, id)
// 	if err != nil {
// 		return err
// 	}
//
// 	item.Children = newChildren
//
// 	return tx.Save(item).Error
// }
//
// // AddFileChildren ... 新增 file 文件树
// func AddFileChildren(tx *gorm.DB, isFatherProject bool, fatherId uint32, file *FileModel) error {
// 	id := file.ID
//
// 	item := GetFolder(fatherId)
//
// 	newChildren, err := addChildren(item.Children, id)
// 	if err != nil {
// 		return err
// 	}
//
// 	item.Children = newChildren
//
// 	return tx.Save(item).Error
// }
//
// type Filer interface {
// 	GetFatherId() uint32
// 	GetProject() uint32
// 	GetID() uint32
// 	AddChildren(*gorm.DB, bool, uint32) error
// }
//
// func Create(db *gorm.DB, f Filer) (uint32, error) {
// 	tx := db.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()
//
// 	if err := tx.Create(f).Error; err != nil {
// 		tx.Rollback()
// 		return 0, err
// 	}
//
// 	isFatherProject := false
// 	fatherId := f.GetFatherId()
// 	if f.GetFatherId() == 0 {
// 		isFatherProject = true
// 		fatherId = f.GetProject()
// 	}
//
// 	if err := f.AddChildren(tx, isFatherProject, fatherId); err != nil {
// 		tx.Rollback()
// 		return 0, err
// 	}
//
// 	return f.GetID(), tx.Commit().Error
// }
