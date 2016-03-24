//Routes for Application
app.config(['$routeProvider', function($routeProvider) {	
	$routeProvider.
	when('/home', {
	   templateUrl: 'app/home/home.view.htm',
	   controller: 'homeController',
	   title : 'Welcome'
	}).
	when('/models', {
	   templateUrl: 'app/models/models.view.htm',
	   controller: 'modelsController',
	   title : 'Models'
	}).
	when('/views', {
	   templateUrl: 'app/views/views.view.htm',
	   controller: 'viewsController',
	   title : 'Views'
	}).
	when('/controllers', {
	   templateUrl: 'app/controllers/controllers.view.htm',
	   controller: 'controllersController',
	   title : 'Controllers'
	}).
	when('/apis', {
	   templateUrl: 'app/apis/apis.view.htm',
	   controller: 'apisController',
	   title : 'APIs'
	}).
	otherwise({
	   redirectTo: '/home'
	});
}]);

