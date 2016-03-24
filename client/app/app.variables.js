var _fbAppId =''
var _gpClientId = ''
//
app.constant('appName','CrowdPuller');
app.constant('appVersion','V1.0');
app.constant('compName','Bitwinger.com');
app.constant('apiPath','https://api.bitwinger.com/crowdpuller');
app.value('appVars', {
	user : {
			sessionId:'',
			memberId:'',
			userName:'',
			fbToken:'',
			gpToken:'',
			address:'',
			isNewSignUp : false
		},
	fbAppId : _fbAppId,
	gpClientId : _gpClientId,
	fbSDKLoaded: false,
	gpSDKLoaded: false,
	fbSDKLoadedHndlr: null,
	gpSDKLoadedHndlr: null
});
