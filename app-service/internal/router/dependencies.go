package router

import (
	"fmt"

	"app-service/internal/config"
	"app-service/internal/controller"
	"app-service/internal/kafka"
	"app-service/internal/repository"
	"app-service/internal/service"

	"go.mongodb.org/mongo-driver/mongo"
)

type Dependencies struct {
	BulletinRepo       repository.IBulletinRepository
	BulletinService    service.IBulletinService
	BulletinController controller.IBulletinController

	GroupRepo       repository.IGroupRepository
	GroupService    service.IGroupService
	GroupController controller.IGroupController

	MemberRepo       repository.IMemberRepository
	MemberService    service.IMemberService
	MemberController controller.IMemberController
}

func NewDependencies(mongoDB *mongo.Database, conf *config.Config) *Dependencies {

	// Initialize repositories
	bulletinRepo := repository.NewBulletinRepository(mongoDB, "bulletins")
	groupRepo := repository.NewGroupRepository(mongoDB, "groups")
	memberRepo := repository.NewMemberRepository(mongoDB, "members")

	// Initialize producers
	notificationProducer, err := kafka.NewNotificationProducer(conf.Kafka.Broker, conf.Kafka.NotificationTopic)
	if err != nil {
		panic(fmt.Sprintf("failed to initialise notification producer: %v", err))
	}

	// Initialize services
	bulletinService := service.NewBulletinService(bulletinRepo)
	groupService := service.NewGroupService(groupRepo, memberRepo, notificationProducer)
	memberService := service.NewMemberService(memberRepo)

	// Initialize controllers
	bulletinController := controller.NewBulletinController(bulletinService)
	groupController := controller.NewGroupController(groupService)
	memberController := controller.NewMemberController(memberService)

	return &Dependencies{
		BulletinRepo:       bulletinRepo,
		BulletinService:    bulletinService,
		BulletinController: bulletinController,

		GroupRepo:       groupRepo,
		GroupService:    groupService,
		GroupController: groupController,

		MemberRepo:       memberRepo,
		MemberService:    memberService,
		MemberController: memberController,
	}
}
