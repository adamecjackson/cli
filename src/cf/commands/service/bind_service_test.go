package service_test

import (
	"cf/api"
	. "cf/commands/service"
	"cf/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testapi "testhelpers/api"
	testassert "testhelpers/assert"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
)

var _ = Describe("bind-service command", func() {
	var (
		requirementsFactory *testreq.FakeReqFactory
	)

	BeforeEach(func() {
		requirementsFactory = &testreq.FakeReqFactory{}
	})

	It("fails requirements when not logged in", func() {
		context := testcmd.NewContext("bind-service", []string{"service", "app"})
		cmd := NewBindService(&testterm.FakeUI{}, testconfig.NewRepository(), &testapi.FakeServiceBindingRepo{})
		testcmd.RunCommand(cmd, context, requirementsFactory)

		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
	})

	Context("when logged in", func() {
		BeforeEach(func() {
			requirementsFactory.LoginSuccess = true
		})

		It("binds a service instance to an app", func() {
			app := models.Application{}
			app.Name = "my-app"
			app.Guid = "my-app-guid"
			serviceInstance := models.ServiceInstance{}
			serviceInstance.Name = "my-service"
			serviceInstance.Guid = "my-service-guid"
			requirementsFactory.Application = app
			requirementsFactory.ServiceInstance = serviceInstance
			serviceBindingRepo := &testapi.FakeServiceBindingRepo{}
			ui := callBindService([]string{"my-app", "my-service"}, requirementsFactory, serviceBindingRepo)

			Expect(requirementsFactory.ApplicationName).To(Equal("my-app"))
			Expect(requirementsFactory.ServiceInstanceName).To(Equal("my-service"))

			testassert.SliceContains(ui.Outputs, testassert.Lines{
				{"Binding service", "my-service", "my-app", "my-org", "my-space", "my-user"},
				{"OK"},
				{"TIP"},
			})
			Expect(serviceBindingRepo.CreateServiceInstanceGuid).To(Equal("my-service-guid"))
			Expect(serviceBindingRepo.CreateApplicationGuid).To(Equal("my-app-guid"))
		})

		It("warns the user when the service instance is already bound to the given app", func() {
			app := models.Application{}
			app.Name = "my-app"
			app.Guid = "my-app-guid"
			serviceInstance := models.ServiceInstance{}
			serviceInstance.Name = "my-service"
			serviceInstance.Guid = "my-service-guid"
			requirementsFactory.Application = app
			requirementsFactory.ServiceInstance = serviceInstance
			serviceBindingRepo := &testapi.FakeServiceBindingRepo{CreateErrorCode: "90003"}
			ui := callBindService([]string{"my-app", "my-service"}, requirementsFactory, serviceBindingRepo)

			testassert.SliceContains(ui.Outputs, testassert.Lines{
				{"Binding service"},
				{"OK"},
				{"my-app", "is already bound", "my-service"},
			})
		})

		It("fails with usage when called without a service instance and app", func() {
			serviceBindingRepo := &testapi.FakeServiceBindingRepo{}

			ui := callBindService([]string{"my-service"}, requirementsFactory, serviceBindingRepo)
			Expect(ui.FailedWithUsage).To(BeTrue())

			ui = callBindService([]string{"my-app"}, requirementsFactory, serviceBindingRepo)
			Expect(ui.FailedWithUsage).To(BeTrue())

			ui = callBindService([]string{"my-app", "my-service"}, requirementsFactory, serviceBindingRepo)
			Expect(ui.FailedWithUsage).To(BeFalse())
		})
	})
})

func callBindService(args []string, requirementsFactory *testreq.FakeReqFactory, serviceBindingRepo api.ServiceBindingRepository) (fakeUI *testterm.FakeUI) {
	fakeUI = new(testterm.FakeUI)
	ctxt := testcmd.NewContext("bind-service", args)

	config := testconfig.NewRepositoryWithDefaults()

	cmd := NewBindService(fakeUI, config, serviceBindingRepo)
	testcmd.RunCommand(cmd, ctxt, requirementsFactory)
	return
}
