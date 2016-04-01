var gapiLoaded = false;
var _fbAppId= '<Enter Your FB AppId>';
var _gpClientId = '<Enter Your Google Client Id>';

function logConsole(msg){
	console.log(msg);
}

function OnGapiLoaded(){
	gapi.load('auth2', function() {//load in the auth2 api's, without it gapi.auth2 will be undefined
		gapi.auth2.init({client_id:_gpClientId});
		gapiLoaded = true;
	});
}

function IsGapiLoaded(){
	return gapiLoaded;
}

//
function getToken(userObj){
	return ((userObj.fbToken != null && userObj.fbToken.length>0)?userObj.fbToken:userObj.gpToken);
}
function checkPageAccess($location, userObj){
	return;//remove this later nimesh
	//if user not logged in then redirect to login page
	if (userObj.memberId == ""){
		$location.path("/login");
	}
	else if (userObj.isNewSignUp == 1){
		$location.path("/profile");
	}
}

function clearAPIError($scope){
	$scope.errors = {api401:false, api404:false, api500:false, apiMsg:null};
}

function handleAPIError($scope, response){
	$scope.errors = {api401:false, api404:false, api500:false, apiMsg:null};
	if (response.status == 404 && response.data.errors){//Validation Error
		if (response.data.errors != null){
			$scope.errors.api404 = true;
			$scope.errors.apiMsg = response.data.errors;
		}
	}
	else if (response.status == 500){//Internal Server Error
		$scope.errors.api500 = true;
		$scope.errors.apiMsg = "We could not process your request due to some problem. Please try again in few minutes.";
	}
	else if (response.status == 401){//Unauthorized Access
		$scope.errors.api401 = true;
		$scope.errors.apiMsg = "You are not logged in or your login session might have timedout. Please copy any changes to clipboard and click here to login.";
	}
}

function ValidateEmail(email){
	pattern=new RegExp("[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,3}$");
	return pattern.test(email);
}

function ValidateUrl(url){
	pattern=new RegExp("https?://.+");
	return pattern.test(url);
}

function IsAlpha(val){
	pattern=new RegExp(/^[a-z]+$/i);
	return pattern.test(val);
}

function IsAlphaNumeric(val){
	pattern=new RegExp(/^[a-z0-9]+$/i);
	return pattern.test(val);
}