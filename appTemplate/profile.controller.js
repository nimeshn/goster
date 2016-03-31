app.controller('profileController', 
	['$scope', '$http', '$filter', '$location', 'apiPath', 'appName', 'appVars', 
		function($scope, $http, $filter, $location, apiPath, appName, appVars) {
	//Shows alert to users
	$scope.showAlert = function(){
		$scope.alertTitle = 'New User SignUp Info';
		$scope.alertMessage = appName + ' requires you to fill in additional details to save & continue using the ' + appName + '.';
		$('#alertModal').modal('show');
	};
	//Get Member Data
	$scope.loadProfile = function(){
		$http.get(apiPath + "/user/" + appVars.user.memberId)
			.then(function(response) {
				if (response.status == 200){
					$scope.profileData = response.data;
					clearAPIError($scope);
				}				
			},
			function(response) {
				handleAPIError($scope, response);
			}
		);
	}	
	//
	$scope.saveProfile =function(){
		$http.put(apiPath + "/user", profileData).
			then(function(response) {
				if (response.status == 200){
					//once profile is saved then it is not 
					appVars.user.isNewSignUp = false;
					clearAPIError($scope);
				}
			  },
			  function(response){
				  handleAPIError($scope, response);
			  });
    };
	//check if the user has access to this page
	checkPageAccess($location, appVars.user);
	if (appVars.user.isNewSignUp){
		$scope.showAlert();		
	}
	//Load Profile
	$scope.loadProfile();	
}]);