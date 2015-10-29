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
}

var IndexViewModel = function() {
  var self = this;
  var timeLastLogin = localStorage.getItem("loginTime");
  this.isLoggedIn = ko.observable(timeLastLogin && moment(timeLastLogin).add(1, 'd').isAfter(moment()))



  this.signOut = function() {
    localStorage.removeItem("token");
    localStorage.removeItem("loginTime");
    this.isLoggedIn(false);
  }

  this.LoginModel = new LoginModel();
  this.RegisterModel = new RegisterModel();

  this.documents = ko.observableArray();

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

  if (this.isLoggedIn()){
    self.LoadData();
  }
}

$(function() {
  $('ul.tabs').tabs();

  var viewModel = new IndexViewModel();
  ko.applyBindings(viewModel);
})
