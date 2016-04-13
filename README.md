# goster

Goster is web framework in written in Golang which will autogenerate AngularJS frontend and Go based Rest API's. Goster enables developers to eaily scaffold a Golang project with 
	1. Define Models for each entity and autogenerate their respective CRUD functions and their wrapper calling method. Currently it supports on MySQL. More databases to come later once the first draft is completed.
	2. Autogenerates Server side Go Controllers for each Model with CRUD functions and Rest API endpoints
	3. Hooks up Rest API's with Golang Http handlers
	4. AutoGenerates angulajs view and controllers for each Model for all CRUD operations.

Each Model contains fields which are statically typed as everything in GoLang is. Each Model consists of Fields which will be mapped to a Database Table and its columns respectively. you can define validations at field level which Goster will use to autogenerate client (html and jquery) and server side validations (Go Controller).

Goster also has capabilities to add social login capabilities to your project. Currently it supports facebook and google account.


# Prerequisites

you would need to get Html prettifier written in golang from https://github.com/yosssi/gohtml.git using the following command
go get github.com/yosssi/gohtml

Also needed is another JS Beautifier "JSBeautifier-go" from https://github.com/ditashi/jsbeautifier-go using the following command
go get github.com/ditashi/jsbeautifier-go