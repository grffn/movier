var LoginModel = function() {
  var self = this;
  this.UserId = ko.observable();
  this.Password = ko.observable();
  this.signIn = function() {
    $.post("/login", JSON.stringify({
      userid: self.UserId(),
      password: self.Password()
    })).success(function(data) {
      localStorage.setItem("token", data.token);
      localStorage.setItem("loginTime", new Date())
    }).error(function(jqxhr, status, error) {
      Materialize.toast("Error: " + error);
    })
  }
}

var RegisterModel = function() {
  var self = this;
  this.Username = ko.observable();
  this.Email = ko.observable();
  this.Password = ko.observable();
  this.register = function() {
    $.post("/register", JSON.stringify({
      username: self.Username(),
      password: self.Password(),
      email: self.Email()
    })).success(function(data) {
      this.Username("");
      this.Email("");
      this.Password("");
    }).error(function(jqxhr, status, error) {
      Materialize.toast("Error: " + error);
    })
  }
};

var DocumentModel = function() {
  var self = this;
  this.Name = ko.observable();
  this.Category = ko.observable();
  this.Tags = ko.observable();
  this.Authors = ko.observable();
  this.file = null;
  this.URL = ko.observable();

  this.fileOpened = function(model, evt) {
    var file = evt.target.files[0];
    self.file = file;
    $.ajax({
      url: "/sign",
      data: {
        filename: file.name,
        filetype: file.type
      },
      method: "get",
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem("token")
      }
    }).success(function(data) {
      var xhr = new XMLHttpRequest();
      xhr.open("PUT", data.signedRequest);
      xhr.setRequestHeader('x-amz-acl', 'public-read');
      xhr.onload = function() {
        if (xhr.status === 200) {
          alert("success")
        }
      };
      xhr.onerror = function() {
        alert("Could not upload file.");
      };
      xhr.send(file);
    }).error(function() {
      console.log(arguments);
    })
  };
}

var IndexViewModel = function() {
  var self = this;
  var timeLastLogin = localStorage.getItem("loginTime");
  this.isLoggedIn = ko.observable(timeLastLogin && moment(timeLastLogin).add(1, 'd').isAfter(moment()))

  this.signOut = function() {
    localStorage.removeItem("token");
    localStorage.removeItem("loginTime");
    this.isLoggedIn(false);
  };

  this.openCreateNew = function(model, evt) {
    $(evt.target).leanModal();
  };

  this.NewDocument = new DocumentModel();

  this.LoginModel = new LoginModel();
  this.RegisterModel = new RegisterModel();

  this.documents = ko.observableArray();

  this.createDocument = function(obj) {
    var data = {
      Name: obj.Name(),
      Category: obj.Category(),
      Tags: obj.Tags().split(/\s*,\s*/),
      Authors: obj.Authors().split(/\s*,\s*/),
      URL: obj.URL()
    }
    if (obj.file) {
      data.MimeType = obj.file.type;
    }
    $.ajax({
      url: "/create",
      method: "post",
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem("token")
      },
      data: JSON.stringify(data)
    }).success(function(){
      self.documents.push(data);
    });
  }

  this.LoadData = function() {
    $.ajax({
      url: "/documents",
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem("token")
      },
      method: 'get'
    }).success(function(data) {
      self.documents(data);
    }).error(function(jqxhr, status, error) {
      Materialize.toast("Error: " + error);
    });
  };

  if (this.isLoggedIn()) {
    self.LoadData();
  }
}

$(function() {
  $('ul.tabs').tabs();

  var viewModel = new IndexViewModel();
  ko.applyBindings(viewModel);
})
