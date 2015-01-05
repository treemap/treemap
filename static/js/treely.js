var mapContainer = L.map('map-container').setView([37.09024, -95.712891], 4);

L.tileLayer('http://server.arcgisonline.com/ArcGIS/rest/services/World_Topo_Map/MapServer/tile/{z}/{y}/{x}', {
    attribution: 'Tiles &copy; Esri &mdash; Esri, DeLorme, NAVTEQ, TomTom, Intermap, iPC, USGS, FAO, NPS, NRCAN, GeoBase, Kadaster NL, Ordnance Survey, Esri Japan, METI, Esri China (Hong Kong), and the GIS User Community'
}).addTo(mapContainer);


var AddGeoJsonsToMap = function(geoms, map) {
    for(var i = 0; i < geoms.length; i++) {
        L.geoJson(JSON.parse(geoms[i])).addTo(map);
    }
}

angular.module('treelyApp', ['ngRoute', 'chieffancypants.loadingBar', 'ngAnimate'])
    .config(function($routeProvider) {
        $routeProvider
            .when('/', {
                redirectTo: "/zipcode/94103"
            })
            .when('/zipcode/:zipcode', {
                controller:'ShowZipcodeCtrl',
                templateUrl:'../templates/zipcodes/show.html'
            })
            .when('/trees', {
                controller:'TreesCtrl',
                templateUrl:'../templates/trees/index.html'
            })
            .when('/trees/:treeId', {
                controller:'ShowTreeCtrl',
                templateUrl:'../templates/trees/show.html'
            })
            .when('/parks/nearby', {
                controller:'NearbyParksCtrl',
                templateUrl:'../templates/parks/index.html'
            })
            .when('/parks', {
                controller:'ParksCtrl',
                templateUrl:'../templates/parks/index.html'
            })
            .when('/lakes', {
                controller:'LakesCtrl',
                templateUrl:'../templates/lakes/index.html'
            })
            .when('/rivers/nearby', {
                controller:'NearbyRiversCtrl',
                templateUrl:'../templates/rivers/index.html'
            })
            .otherwise({
                redirectTo:'/'
            });
    })
    .config(function(cfpLoadingBarProvider) {
        cfpLoadingBarProvider.includeSpinner = true;
    })
    .controller('NearbyParksCtrl', function($scope, $http, cfpLoadingBar) {
        $scope.parks = {}

        $scope.$watch("zipcode", function(zipcode, oldValue) {
            $http.get(SarpaServiceDiscovery.treemap[0] + '/v1/zipcodes/' + zipcode.number + '/parks').
                success(function(data, status, headers, config) {
                    cfpLoadingBar.start();
                    $scope.parks = data;

                    for(var i = 0; i < $scope.parks.length; i++) {
                        cfpLoadingBar.inc();

                        console.log(i);
                        L.geoJson(JSON.parse($scope.parks[i].geom)).addTo(mapContainer);
                    }
                    cfpLoadingBar.complete()
                }).
                error(function(data, status, headers, config) {});
        });
    })
    .controller('NearbyTreesCtrl', function($scope, $http, cfpLoadingBar) {
        $scope.trees = [];

        $scope.$watch("zipcode", function(zipcode, oldValue) {
            $http.get(SarpaServiceDiscovery.treemap[0] + '/v1/zipcodes/' + zipcode.number + '/trees').
                success(function(data, status, headers, config) {
                    $scope.trees = data;
                }).
                error(function(data, status, headers, config) {
                });
        })
    })
    .controller('NearbyHydrologyCtrl', function($scope, $http, cfpLoadingBar) {
        $scope.hydrology = {};
        $scope.hydroType = "lakes";

        $scope.init = function(hydroType) {
            $scope.hydroType = hydroType;
        }

        $scope.$watch("zipcode", function(zipcode, oldValue) {
            $http.get(SarpaServiceDiscovery.treemap[0] + '/v1/zipcodes/' + zipcode.number + '/' + $scope.hydroType).
                success(function(data, status, headers, config) {
                    cfpLoadingBar.start();
                    $scope.hydrology = data;

                    for(var i = 0; i < $scope.hydrology.length; i++) {
                        cfpLoadingBar.inc();

                        console.log(i);
                        L.geoJson(JSON.parse($scope.hydrology[i].geom)).addTo(mapContainer);
                    }
                    cfpLoadingBar.complete()

                }).
                error(function(data, status, headers, config) {});
        });
    })
    .controller('ShowTreeCtrl', function($scope, $http, $routeParams) {
        $scope.tree = {}

        $http.get('/data/trees/' + $routeParams.treeId + '.json').
            success(function(data, status, headers, config) {
                $scope.tree = data;

                AddGeoJsonsToMap($scope.tree.geom, mapContainer);

                var center = JSON.parse($scope.tree.center);
                mapContainer.setView(center.coordinates.reverse(), 6);
            }).
            error(function(data, status, headers, config) {});
    })
    .controller('TreesCtrl', function($scope, $http) {
        $scope.trees = [];

        $http.get('/data/trees/index.json').
            success(function(data, status, headers, config) {
                $scope.trees = data;
            }).
            error(function(data, status, headers, config) {});

    })
    .controller('ParksCtrl', function($scope, $http, $routeParams, cfpLoadingBar) {
        $scope.parks = {}

        $http.get('/data/parks/index.json').
            success(function(data, status, headers, config) {
                cfpLoadingBar.start();
                $scope.parks = data;

                for(var i = 0; i < $scope.parks.length; i++) {
                    cfpLoadingBar.inc();

                    L.geoJson(JSON.parse($scope.parks[i].geom)).addTo(mapContainer);
                }
                cfpLoadingBar.complete()

            }).
            error(function(data, status, headers, config) {});
    })
    .controller('LakesCtrl', function($scope, $http, $routeParams, cfpLoadingBar) {
        $scope.lakes = {}

        $http.get(SarpaServiceDiscovery.treemap[0] + '/lakes').
            success(function(data, status, headers, config) {
                cfpLoadingBar.start();
                $scope.lakes = data;

                for(var i = 0; i < $scope.lakes.length; i++) {
                    cfpLoadingBar.inc();

                    L.geoJson(JSON.parse($scope.lakes[i].geom)).addTo(mapContainer);
                }
                cfpLoadingBar.complete()

            }).
            error(function(data, status, headers, config) {});
    })
    .controller('RiversCtrl', function($scope, $http, $routeParams, cfpLoadingBar) {
        $scope.rivers = {}

        $http.get(SarpaServiceDiscovery.treemap[0] + '/rivers').
            success(function(data, status, headers, config) {
                cfpLoadingBar.start();
                $scope.rivers = data;

                for(var i = 0; i < $scope.rivers.length; i++) {
                    cfpLoadingBar.inc();

                    L.geoJson(JSON.parse($scope.rivers[i].geom)).addTo(mapContainer);
                }
                cfpLoadingBar.complete()

            }).
            error(function(data, status, headers, config) {});
    })
    .controller('ShowZipcodeCtrl', function($scope, $http, $routeParams) {
        $scope.zipcode = {}

        $http.get('/data/zipcodes/' + $routeParams.zipcode + '.json').
            success(function(data, status, headers, config) {
                $scope.zipcode = data;

                // Shouldn't need to parse this as well. Should be sent parsed.
                L.geoJson(JSON.parse($scope.zipcode.geom)).addTo(mapContainer);

                var center = JSON.parse($scope.zipcode.center);
                mapContainer.setView(center.coordinates.reverse(), 10);
            }).
            error(function(data, status, headers, config) {});
    });
