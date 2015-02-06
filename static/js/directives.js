/**
 * NEUBOARD - Responsive Admin Theme
 * Copyright 2014 Authentic Goods Co. http://authenticgoods.co
 *
 * TABLE OF CONTENTS
 * Use @ along with function name to search for the directive.
 *
 *  @pageTitle - Page Title Directive for page title name
 *  @widgetToggle - Directive to toggle widgets
 *  @widgetClose - Directive to close widget
 *  @toggleLeftSidebar - Left Sidebar Directive to toggle sidebar navigation
 *  @toggleProfile - Show/Hide Profile View
 *  @toggleRightSidebar - Right Sidebar Directive to toggle sidebar navigation
 *  @navToggleSub - Directive to toggle sub-menu down
 *  @slider - Directive to run bootstrap sliders
 *  @gaugejs - Directive for the gauge graph
 *  @css3animate - css3 animations
 *  @mapUsa - Directive for vector map
 *  @iCheck - Directive for custom checkboxes
 *  @skycon - Directive for skycons
 *  @dropZone - Directive for multifile uploader
 *  @c3chart - Directive for c3chart
 *  @chartColumn - Directive for c3chart/chartColumn
 *  @chartAxes - Directive for c3chart/chartAxes
 *  @chartAxis - Directive for c3chart/chartAxis
 *  @chartAxisX - Directive for c3chart/chartAxisX
 *  @chartAxisY - Directive for c3chart/chartAxisY
 *  @chartGrid - Directive for c3chart/chartGrid
 *  @chartGridOptional - Directive for c3chart/chartGridOptional
 *  @chartAxisXTick - Directive for c3chart/chartAxisXTick
 *  @chartLegend - Directive for c3chart/chartLegend
 *  @chartTooltip - Directive for c3chart/chartTooltip
 *  @chartSize - Directive for c3chart/chartSize
 *  @chartColors - Directive for c3chart/chartColors
 *  @barchart - Directive for morris/barchart
 *  @linechart - Directive for morris/linechart
 *  @donutchart - Directive for morris/donutchart
 *  @sparkline - Directive for sparkline chart
 *  @toggleSettings - Directive to toggle settings widgets for DEMO
 *  @switchTheme - Directive to switch theme colors for DEMO
 *
 */

/*
 * @pageTitle - Page Title Directive for page title name
 */
function pageTitle($rootScope, $timeout) {
    return {
        link: function(scope, element) {
            var listener = function(event, toState, toParams, fromState, fromParams) {
                var title = 'NeuBoard - Responsive Admin Theme';
                if (toState.data && toState.data.pageTitle) title = 'NeuBoard | ' + toState.data.pageTitle;
                $timeout(function() {
                    element.text(title);
                });
            };
            $rootScope.$on('$stateChangeStart', listener);
        }
    }
};

/**
 * @widgetToggle - Directive to toggle widget
 */
function widgetToggle() {
    return {
        restrict: 'A',
        link: function(scope, element) {
            element.click(function() {
                $(this).parent().parent().next().slideToggle("fast"), $(this).toggleClass("fa-chevron-down fa-chevron-up")
            });
        }
    }
};

/**
 * @widgetClose - Directive to close widget
 */
function widgetClose() {
    return {
        restrict: 'A',
        link: function(scope, element) {
            element.click(function() {
                $(this).parent().parent().parent().fadeOut()
            });
        }
    }
};

/*
 * @toggleLeftSidebar - Left Sidebar Directive to toggle sidebar navigation
 */
function toggleLeftSidebar() {
    return {
        restrict: 'A',
        template: '<button ng-click="toggleLeft()" class="sidebar-toggle" id="toggle-left"><i class="fa fa-bars"></i></button>',
        controller: function($scope, $element) {
            $scope.toggleLeft = function() {
                ($(window).width() > 767) ? $('#main-wrapper').toggleClass('sidebar-mini'): $('#main-wrapper').toggleClass('sidebar-opened');
            }
        }
    };
}

/*
 * @toggleProfile - Show/Hide Profile View
 */
function toggleProfile() {
    return {
        restrict: 'A',
        template: '<button ng-click="toggleProfile()" type="button" class="btn btn-default" id="toggle-profile"><i class="icon-user"></i></button>',
        controller: function($scope, $element) {
            $scope.toggleProfile = function() {
                $('.sidebar-profile').slideToggle();
            }
        }
    };
};

/*
 * @toggleRightSidebar - Right Sidebar Directive to toggle sidebar navigation
 */
function toggleRightSidebar() {
    return {
        restrict: 'A',
        template: '<button ng-click="toggleRight()" class="sidebar-toggle" id="toggle-right"><i class="fa fa-indent"></i></button>',
        controller: function($scope, $element) {
            $scope.toggleRight = function() {
                $('#sidebar-right').toggleClass("sidebar-right-open");
                $("#toggle-right .fa").toggleClass("fa-indent fa-dedent");
            }
        }
    };
};

/**
 * @navToggleSub - Directive to toggle sub-menu down
 */
function navToggleSub() {
    return {
        restrict: 'A',
        link: function(scope, element) {
            element.navgoco({
                caretHtml: false,
                accordion: true
            })
        }
    };
};

/**
 * @sliders - Directive to run bootstrap sliders
 */
function slider() {
    return {
        restrict: 'A',
        link: function(scope, element) {
            element.slider();
        }
    }
};

/**
 * @gaugejs - Directive for the gauge graph
 */
function gauge() {
    return {
        restrict: 'AC',
        scope: {
            'animationTime': '=',
            'value': '=',
            'options': '=',
            'maxValue': '=',
            'gaugeType': '='
        },
        controller: function($scope, $element) {
            if ($scope.gaugeType === 'donut') {
                $scope.gauge = new Donut($element[0]);
            } else {
                $scope.gauge = new Gauge($element[0]);
            }
            $scope.gauge.maxValue = $scope.maxValue;
            $scope.$watchCollection('[options, value]', function(newValues) {
                $scope.gauge.setOptions(newValues[0]);
                if (!isNaN(newValues[1])) {
                    $scope.gauge.set(newValues[1]);
                }
            });
        },
    }
};


/**
 * @css3animate - css3 animations
 */
function css3animate() {
    return {
        restrict: 'A',
        link: function(scope, element) {
            element.click(function() {
                var animate = $(this).attr("data-animated");
                $(this).closest('.panel').addClass(animate).delay(1000).queue(function(next) {
                    $(this).removeClass(animate);
                    next();
                })
            })
        }
    }

};


/**
 * @mapUsa - Directive for vector map
 */
function mapUsa() {
    return {

        restrict: 'A',
        link: function(scope, element, attrs) {
            element.vectorMap({
                map: 'us_aea_en',
                backgroundColor: 'transparent',
                zoomButtons: true,
                regionStyle: {
                    initial: {
                        fill: '#909AA0'
                    },
                    hover: {
                        fill: '#1D212A'
                    }
                },
                onRegionClick: function(event, code) {

                },
                markerStyle: {
                    initial: {
                        fill: '#27B6AF',
                        stroke: '#27B6AF',
                    }
                },
                markers: [{
                    latLng: [37.78, -122.41],
                    name: 'San Francisco',
                    style: {
                        r: 10
                    }
                }, {
                    latLng: [40.71, -74],
                    name: 'New York City',
                    style: {
                        r: 15
                    }
                }, {
                    latLng: [41.89, -87.62],
                    name: 'Chicago',
                    style: {
                        r: 5
                    }
                }, {
                    latLng: [34.00, -118.25],
                    name: 'Los Angeles',
                    style: {
                        r: 20
                    }
                }, {
                    latLng: [34.00, -106.00],
                    name: 'New Mexico',
                    style: {
                        r: 10
                    }
                }, {
                    latLng: [44.50, -100.00],
                    name: 'South Dakota',
                    style: {
                        r: 13
                    }
                }, {
                    latLng: [25.78, -80.22],
                    name: 'Miami',
                    style: {
                        r: 7
                    }
                }, ]
            });

        }

    }

};

/**
 * @iCheck - Directive for custom checkboxes
 */
function ichecks($timeout, $parse) {
    return {
        link: function(scope, element, attrs) {
            return $timeout(function() {
                return $(element).iCheck({
                    checkboxClass: 'icheckbox_flat-grey',
                    radioClass: 'iradio_flat-grey',
                    increaseArea: '20%'
                });
            });
        }
    };
};

/**
 * @skycon - Directive for skycons
 */
function skycon() {
    return {
        restrict: 'A',
        link: function(scope, element, attrs) {
            var skycons = new Skycons({
                'color': (attrs.color || '#27B6AF')
            });
            element.html('<canvas width="' + attrs.width + '" height="' + attrs.height + '"></canvas>');
            skycons.add(element.children()[0], attrs.skycon);
            skycons.play()
        }
    }
};

/**
 * Directive for the calendar widget
 */

(function() {
    var TienClndrDirective, module;

    module = angular.module('tien.clndr', []);

    TienClndrDirective = function() {
        var controller, scope;
        scope = {
            clndr: '=tienClndrObject',
            events: '=tienClndrEvents',
            options: '=?tienClndrOptions'
        };
        controller = function($scope, $element, $attrs, $transclude) {
            return $transclude(function(clone, scope) {
                var options, render;
                $element.append(clone);
                $scope.$watch('events', function(val) {
                    if (val != null ? val.length : void 0) {
                        return $scope.clndr.setEvents(angular.copy(val));
                    }
                });
                render = function(data) {
                    return angular.extend(scope, data);
                };
                options = angular.extend($scope.options || {}, {
                    render: render
                });
                return $scope.clndr = angular.element("<div/>").clndr(options);
            });
        };
        return {
            restrict: 'E',
            replace: true,
            transclude: true,
            scope: scope,
            controller: controller
        };
    };

    module.directive('tienClndr', TienClndrDirective);

}).call(this);

/**
 * dropZone - Directive for multifile uploader
 */
function dropZone() {
    return function(scope, element, attrs) {
        element.dropzone({
            url: "/upload",
            maxFilesize: 100,
            paramName: "uploadfile",
            maxThumbnailFilesize: 5
        });
    }
}


/**
 * c3chart - Directive for c3chart
 */
function c3chart($timeout) {
    var chartLinker = function(scope, element, attrs, chartCtrl) {
        // Trick to wait for all rendering of the DOM to be finished.
        $timeout(function() {
            chartCtrl.showGraph()
        });
    };

    return {
        "restrict": "E",
        "controller": "ChartController",
        "scope": {
            "bindto": "@bindtoId",
            "showLabels": "@showLabels",
            "showSubchart": "@showSubchart",
            "enableZoom": "@enableZoom",
            "chartData": "=chartData",
            "chartColumns": "=chartColumns",
            "chartX": "=chartX"
        },
        "template": "<div><div id='{{bindto}}'></div><div ng-transclude></div></div>",
        "replace": true,
        "transclude": true,
        "link": chartLinker
    }
};

/**
 * chartColumn - Directive for c3chart/chartColumn
 */
function chartColumn() {
    var columnLinker = function(scope, element, attrs, chartCtrl) {
        var column = attrs['columnValues'].split(",");
        column.unshift(attrs['columnId']);
        chartCtrl.addColumn(column, attrs['columnType'], attrs['columnName'], attrs['columnColor']);
    };

    return {
        "require": "^c3chart",
        "restrict": "E",
        "scope": {},
        "replace": true,
        "link": columnLinker
    }
};

/**
 * chartAxes - Directive for c3chart/chartAxes
 */
function chartAxes() {
    var axesLinker = function(scope, element, attrs, chartCtrl) {
        var x = attrs['valuesX'];
        if (x) {
            chartCtrl.addXAxisValues(x);
        }

        var y = attrs['y'];
        var y2 = attrs['y2'];
        var yAxis = {};
        if (y2) {
            var items = y2.split(",");
            for (item in items) {
                yAxis[items[item]] = "y2";
            }
            if (y) {
                var items = y.split(",");
                for (item in items) {
                    yAxis[items[item]] = "y";
                }
            }
            chartCtrl.addYAxis(yAxis);
        }
    };

    return {
        "require": "^c3chart",
        "restrict": "E",
        "scope": {},
        "replace": true,
        "link": axesLinker
    }

};

/**
 * chartAxis - Directive for c3chart/chartAxis
 */
function chartAxis() {
    var axisLinker = function(scope, element, attrs, chartCtrl) {
        var rotate = attrs['axisRotate'];
        if (rotate) {
            chartCtrl.rotateAxis();
        }
    };

    return {
        "require": "^c3chart",
        "restrict": "E",
        "scope": {},
        "transclude": true,
        "template": "<div ng-transclude></div>",
        "replace": true,
        "link": axisLinker
    }

};

/**
 * chartAxisX - Directive for c3chart/chartAxisX
 */
function chartAxisX() {
    var axisLinker = function(scope, element, attrs, chartCtrl) {
        var position = attrs['axisPosition'];
        var label = attrs['axisLabel'];

        var axis = {
            "label": {
                "text": label,
                "position": position
            }
        };

        var type = attrs['axisType'];
        if (type) {
            axis.type = type;
        }

        var height = attrs['axisHeight'];
        if (height) {
            axis.height = height;
        }
        chartCtrl.addAxisProperties('x', axis);
    };

    return {
        "require": "^c3chart",
        "restrict": "E",
        "scope": {},
        "transclude": true,
        "template": "<div ng-transclude></div>",
        "replace": true,
        "link": axisLinker
    }

};

/**
 * chartAxisY - Directive for c3chart/chartAxisY
 */
function chartAxisY() {
    var axisLinker = function(scope, element, attrs, chartCtrl) {
        var id = attrs['axisId'];
        var position = attrs['axisPosition'];
        var label = attrs['axisLabel'];

        var axis = {
            "label": {
                "text": label,
                "position": position
            }
        };
        if (id === 'y2') {
            axis.show = true;
        }
        var paddingTop = attrs['paddingTop'];
        var paddingBottom = attrs['paddingBottom'];
        if (paddingTop | paddingBottom) {
            paddingTop = (paddingTop) ? paddingTop : 0;
            paddingBottom = (paddingBottom) ? paddingBottom : 0;
            axis.padding = {
                "top": parseInt(paddingTop),
                "bottom": parseInt(paddingBottom)
            };
        }
        var rangeMax = attrs['rangeMax'];
        var rangeMin = attrs['rangeMin'];
        if (rangeMax) {
            axis.max = parseInt(rangeMax);
        }
        if (rangeMin) {
            axis.min = parseInt(rangeMin);
        }

        chartCtrl.addAxisProperties(id, axis);
    };

    return {
        "require": "^c3chart",
        "restrict": "E",
        "scope": {},
        "replace": true,
        "link": axisLinker
    }

};

/**
 * chartGrid - Directive for c3chart/chartGrid
 */
function chartGrid() {
    var gridLinker = function(scope, element, attrs, chartCtrl) {
        var showX = attrs["showX"];
        if (showX && showX === "true") {
            chartCtrl.addGrid("x");
        }
        var showY = attrs["showY"];
        if (showY && showY === "true") {
            chartCtrl.addGrid("y");
        }
        var showY2 = attrs["showY2"];
        if (showY2 && showY2 === "true") {
            chartCtrl.addGrid("y2");
        }
    };

    return {
        "require": "^c3chart",
        "restrict": "E",
        "scope": {},
        "replace": true,
        "link": gridLinker,
        "transclude": true,
        "template": "<div ng-transclude></div>"
    }

};


/**
 * chartGridOptional - Directive for c3chart/chartGridOptional
 */
function chartGridOptional() {
    var gridLinker = function(scope, element, attrs, chartCtrl) {
        var axisId = attrs["axisId"];
        var value = attrs["gridValue"];
        var text = attrs["gridText"];

        chartCtrl.addGridLine(axisId, value, text);
    };

    return {
        "require": "^c3chart",
        "restrict": "E",
        "scope": {},
        "replace": true,
        "link": gridLinker
    }

};

/**
 * chartAxisXTick - Directive for c3chart/chartAxisXTick
 */
function chartAxisXTick() {
    var tickLinker = function(scope, element, attrs, chartCtrl) {
        var tick = {};

        var count = attrs['tickCount'];
        if (count) {
            tick.count = count;
        }

        // TODO, dit lijkt nog niet echt iets te doen
        var format = attrs['tickFormat'];
        if (format) {
            tick.format = format;
        }

        var culling = attrs['tickCulling'];
        if (culling) {
            tick.culling = culling;
        }

        var rotate = attrs['tickRotate'];
        if (rotate) {
            tick.rotate = rotate;
        }

        var fit = attrs['tickFit'];
        if (fit) {
            tick.fit = fit;
        }

        chartCtrl.addXTick(tick);
    };

    return {
        "require": "^c3chart",
        "restrict": "E",
        "scope": {},
        "replace": true,
        "link": tickLinker
    }
};


/**
 * chartLegend - Directive for c3chart/chartLegend
 */
function chartLegend() {
    var legendLinker = function(scope, element, attrs, chartCtrl) {
        var legend = null;
        var show = attrs["showLegend"];
        if (show && show === "false") {
            legend = {
                "show": false
            };
        } else {
            var position = attrs["legendPosition"];
            if (position) {
                legend = {
                    "position": position
                };
            }
        }

        if (legend != null) {
            chartCtrl.addLegend(legend);
        }
    };

    return {
        "require": "^c3chart",
        "restrict": "E",
        "scope": {},
        "replace": true,
        "link": legendLinker
    }

};


/**
 * chartTooltip - Directive for c3chart/chartTooltip
 */
function chartTooltip() {
    var tooltipLinker = function(scope, element, attrs, chartCtrl) {
        var tooltip = null;
        var show = attrs["showTooltip"];
        if (show && show === "false") {
            tooltip = {
                "show": false
            };
        } else {
            var grouped = attrs["groupTooltip"];
            if (grouped && grouped === "false") {
                tooltip = {
                    "grouped": false
                };
            }
        }

        if (tooltip != null) {
            chartCtrl.addTooltip(tooltip);
        }
    };

    return {
        "require": "^c3chart",
        "restrict": "E",
        "scope": {},
        "replace": true,
        "link": tooltipLinker
    }

};


/**
 * chartSize - Directive for c3chart/chartSize
 */
function chartSize() {
    var sizeLinker = function(scope, element, attrs, chartCtrl) {
        var chartSize = null;
        var width = attrs["chartWidth"];
        var height = attrs["chartHeight"]
        if (width || height) {
            chartSize = {};
            if (width) {
                chartSize.width = parseInt(width);
            }
            if (height) {
                chartSize.height = parseInt(height);
            }
            chartCtrl.addSize(chartSize);
        }
    };

    return {
        "require": "^c3chart",
        "restrict": "E",
        "scope": {},
        "replace": true,
        "link": sizeLinker
    }

};

/**
 * chartColors - Directive for c3chart/chartColors
 */
function chartColors() {
    var colorsLinker = function(scope, element, attrs, chartCtrl) {
        var pattern = attrs["colorPattern"];
        if (pattern) {
            chartCtrl.addColors(pattern.split(","));
        }
    };

    return {
        "require": "^c3chart",
        "restrict": "E",
        "scope": {},
        "replace": true,
        "link": colorsLinker
    }
};

/**
 * barchart - Directive for morris/barchart
 */
function barchart() {
    function createChart(el_id, options) {
        options.element = el_id;
        var r = new Morris.Bar(options);
        return r;
    }

    return {
        restrict: 'E',
        scope: {
            options: '='
        },
        replace: true,
        template: '<div></div>',
        link: function link(scope, element, attrs) {
            return createChart(attrs.id, scope.options);
        }
    };
};

/**
 * linechart - Directive for morris/linechart
 */
function linechart() {
    function createChart(el_id, options) {
        options.element = el_id;
        var r = new Morris.Line(options);
        return r;
    }

    return {
        restrict: 'E',
        scope: {
            options: '='
        },
        replace: true,
        template: '<div></div>',
        link: function(scope, element, attrs) {
            return createChart(attrs.id, scope.options)
        }
    }
};

/**
 * donutchart - Directive for morris/donutchart
 */
function donutchart() {
    function createChart(el_id, options) {
        options.element = el_id;
        var r = new Morris.Donut(options);
        return r;
    }

    return {
        restrict: 'E',
        scope: {
            options: '='
        },
        replace: true,
        template: '<div></div>',
        link: function(scope, element, attrs) {
            return createChart(attrs.id, scope.options)
        }
    }
};


/**
 * sparkline - Directive for sparkline chart
 */
function sparkline() {
    return {
        restrict: 'A',
        scope: {
            sparkData: '=',
            sparkOptions: '=',
        },
        link: function(scope, element, attrs) {
            scope.$watch(scope.sparkData, function() {
                render();
            });
            scope.$watch(scope.sparkOptions, function() {
                render();
            });
            var render = function() {
                $(element).sparkline(scope.sparkData, scope.sparkOptions);
            };
        }
    }
};
/**
 * @toggleSettings - Directive to toggle settings widgets for DEMO
 */
function toggleSettings() {
    return {
        restrict: 'A',
        link: function(scope, element) {
            element.click(function() {
                if ($(this).hasClass('open')) {
                    $('#config').animate({
                        "right": "-205px"
                    }, 150);
                    $(this).removeClass('open').addClass('closed');
                } else {
                    $("#config").animate({
                        "right": "0px"
                    }, 150);
                    $(this).removeClass('closed').addClass('open');
                }
            });
        }
    }
};

/**
 * @switchTheme - Directive to switch theme colors for DEMO
 */
function switchTheme() {
    return {
        restrict: 'A',
        link: function(scope, element) {
            element.click(function() {
                $('#main-wrapper').attr('class', '');
                var themeValue = $(this).data('theme');
                $('#main-wrapper').addClass(themeValue);
            });
        }
    }
};

/*
 * Pass functions to module
 */
angular
    .module('neuboard')
    .directive('pageTitle', pageTitle)
    .directive('widgetToggle', widgetToggle)
    .directive('widgetClose', widgetClose)
    .directive('toggleLeftSidebar', toggleLeftSidebar)
    .directive('toggleProfile', toggleProfile)
    .directive('toggleRightSidebar', toggleRightSidebar)
    .directive('navToggleSub', navToggleSub)
    .directive('slider', slider)
    .directive('gauge', gauge)
    .directive('css3animate', css3animate)
    .directive('mapUsa', mapUsa)
    .directive('ichecks', ichecks)
    .directive('skycon', skycon)
    .directive('dropZone', dropZone)
    .directive('c3chart', c3chart)
    .directive('chartColumn', chartColumn)
    .directive('chartAxes', chartAxes)
    .directive('chartAxis', chartAxis)
    .directive('chartAxisX', chartAxisX)
    .directive('chartAxisY', chartAxisY)
    .directive('chartGrid', chartGrid)
    .directive('chartGridOptional', chartGridOptional)
    .directive('chartAxisXTick', chartAxisXTick)
    .directive('chartLegend', chartLegend)
    .directive('chartTooltip', chartTooltip)
    .directive('chartSize', chartSize)
    .directive('chartColors', chartColors)
    .directive('barchart', barchart)
    .directive('linechart', linechart)
    .directive('donutchart', donutchart)
    .directive('sparkline', sparkline)
    .directive('toggleSettings', toggleSettings)
    .directive('switchTheme', switchTheme)
