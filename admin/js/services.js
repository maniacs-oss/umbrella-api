'use strict';

/* Services */

var secFirstServices = angular.module('secFirstServices', []);

secFirstServices.factory('Segments', ['$http', '$cookieStore',
  function($http, $cookieStore){
    var factory = {};

    factory.getAll = function() {
        return $http.get('/v1/segments', {
            headers: {'token': $cookieStore.get('token')}
        });
    };

    factory.getRaw = function() {
        return $http.get('/v1/segments_raw', {
            headers: {'token': $cookieStore.get('token')}
        });
    };

    factory.getByCat = function(categoryId) {
        return $http.get("/v1/segments/"+categoryId+"/category");
    };

    factory.getId = function(segmentId) {
        console.log(segmentId);
        return $http.get("/v1/segments_raw/"+segmentId);
    };
    return factory;
  }]);

secFirstServices.factory('Login', ['$http', '$cookieStore',
  function($http, $cookieStore){
    var factory = {};

    factory.postForm = function(formData) {
        return $http({
          method  : 'POST',
          url     : '/v1/account/login',
          headers: {'token': $cookieStore.get('token')},
          data    : formData
      });
    };
    return factory;
  }]);

secFirstServices.factory('Categories', ['$http', '$cookieStore',
  function($http, $cookieStore){
    var factory = {};

    factory.postForm = function(formData) {
        return $http({
          method  : 'POST',
          url     : '/v1/account/login',
          headers: {'token': $cookieStore.get('token')},
          data    : formData
      });
    };
    return factory;
  }]);