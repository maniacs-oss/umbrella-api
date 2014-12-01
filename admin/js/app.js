'use strict';

/* App Module */

var secFirstApp = angular.module('secFirstApp', [
  'ngRoute',
  'ngCookies',
  'secFirstAnimations',
  'secFirstControllers',
  'secFirstFilters',
  'secFirstServices',
  'summernote'
]);

secFirstApp.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/segments', {
        templateUrl: 'partials/segments.html',
        controller: 'SegmentList'
      }).
      when('/login', {
        templateUrl: 'partials/login.html',
        controller: 'LoginForm'
      }).
      when('/about', {
        templateUrl: 'partials/about.html'
      }).
      when('/segments/:segmentId', {
        templateUrl: 'partials/segment-detail.html',
        controller: 'SegmentDetail'
      }).
      otherwise({
        redirectTo: '/segments'
      });
  }]);
