/**
 * @ngdoc directive
 * @name ng.directive:rcSubmit
 *
 * @description
 * Alternative to ngSubmit that verifies the ngFormController is valid before
 * executing the given expression.  Otherwise it cancels the event. 
 *
 * @element form
 * @param {expression} rcSubmit {@link guide/expression Expression} to eval.
 */
var rcSubmitDirective = {
  'rcSubmit': ['$parse', '$q', '$timeout', function ($parse, $q, $timeout) {
    return {
      restrict: 'A',
      require: ['rcSubmit', '?form'],
      controller: ['$scope', function ($scope) {
 
 		var formElement = null;
        var formController = null;
        var attemptHandlers = [];
        var submitCompleteHandlers = [];
 
        this.attempted = false;
        this.submitInProgress = false;
        
        this.setFormElement = function(element) {
          formElement = element;
        }
        
        this.submit = function() {
          if (!formElement) return;
          
          jQuery(formElement).submit();
        }
        
        this.onAttempt = function(handler) {
          attemptHandlers.push(handler);
        };
 
        this.setAttempted = function() {
          this.attempted = true;
          
          angular.forEach(attemptHandlers, function (handler) {
            handler();
          });
        };
 
        this.setFormController = function(controller) {          
          formController = controller;
        };
 
        this.needsAttention = function (fieldModelController) {
          if (!formController) return false;
 
          if (fieldModelController) {
            return fieldModelController.$invalid && 
                   (fieldModelController.$dirty || this.attempted);
          } else {
            return formController && formController.$invalid && 
                   (formController.$dirty || this.attempted);
          }
        };
 
        this.onSubmitComplete = function (handler) {
 
          submitCompleteHandlers.push(handler);
        };
 
        this.setSubmitComplete = function (success, data) {
 
          angular.forEach(submitCompleteHandlers, function (handler) {
            handler({ 'success': success, 'data': data });
          });
        };
      }],
      compile: function(cElement, cAttributes, transclude) {
        return {
          pre: function(scope, formElement, attributes, controllers) {
 
            var submitController = controllers[0];
            var formController = (controllers.length > 1) ? controllers[1] : null;
 
 			submitController.setFormElement(formElement);
            submitController.setFormController(formController);
 
            scope.rc = scope.rc || {};
            scope.rc[attributes.name] = submitController;
          },
          post: function(scope, formElement, attributes, controllers) {
 
            var submitController = controllers[0];
            var formController = (controllers.length > 1) ? controllers[1] : null;
            var fn = $parse(attributes.rcSubmit);
 
            formElement.bind('submit', function (event) {
              submitController.setAttempted();
              if (!scope.$$phase) scope.$apply();
 
              if (!formController.$valid) return false;
 
              var doSubmit = function () {
 
                submitController.submitInProgress = true;
                if (!scope.$$phase) scope.$apply();
 
                var returnPromise = $q.when(fn(scope, { $event: event }));
 
                returnPromise.then(function (result) {
                  submitController.submitInProgress = false;
                  if (!scope.$$phase) scope.$apply();
                  
                  // This is a small hack.  We want the submitInProgress
                  // flag to be applied to the scope before we actually
                  // raise the submitComplete event. We do that by
                  // using angular's $timeout service which even without
                  // a timeout value specified will not fire until after
                  // the scope is digested.
                  $timeout(function() {
                    submitController.setSubmitComplete(true, result);
                  });
 
                }, function (error) {
                  submitController.submitInProgress = false;
                  if (!scope.$$phase) scope.$apply();
                  $timeout(function() {
                    submitController.setSubmitComplete(false, error);
                  });
                });
              };
 
              if (!scope.$$phase) {
                scope.$apply(doSubmit);
              } else {
                doSubmit();
                if (!scope.$$phase) scope.$apply();
              }
            });
          }
        };
      }
    };
  }]
};


/**
 * @ngdoc module
 * @name rcForm
 *
 * @description
 * Module to encapsulate all of our custom form-based directives.
 */
var rcFormModule = angular.module('rcForm', []);

var rcSubmitDirective = rcSubmitDirective || null;
var rcVerifySetDirective = rcVerifySetDirective || null;
if (rcSubmitDirective) rcFormModule.directive(rcSubmitDirective);
if (rcVerifySetDirective) rcFormModule.directive(rcVerifySetDirective);


/**
 * @ngdoc module
 * @name rcDisabled
 *
 * @description
 * Module to encapsulate the rcDisabled directive and rcDisabledProvider.
 */

/**
 * @ngdoc directive
 * @name ng.directive:rcDisabled
 *
 * @description
 * calls rcDisabledProvider.disable on the given element.
 *
 * @element ANY
 * @param {expression} rcDisabled {@link guide/expression Expression} to watch 
 * which determines when to disable.
 */
var rcDisabledDirective = {
  'rcDisabled': ['rcDisabled', function (rcDisabled) {
    return {
      restrict: 'A',
      link: function (scope, element, attributes) {
        
        scope.$watch(attributes.rcDisabled, function(isDisabled) {
          rcDisabled.disable(element, isDisabled);
        });
      }
    }
  }]
};

/**
 * @ngdoc provider
 * @name ng.provider:rcDisabledProvider
 *
 * @description
 * The provider for rcDisabled. Allows configuration of the method used when
 * toggling disabled.
 *
 */
var rcDisabledProvider = function () {
    
  var defaultDisableHandler = function(rootElement, isDisabled) {
    var jElement = jQuery(rootElement);
    
    return jElement
            .find(':not([rc-disabled])')
            .filter(function(index) {
              return jQuery(this)
                       .parents()
                       .not(jElement)
                       .filter('[rc-disabled]').length === 0;
            })
            .filter('input:not([ng-disabled]), button:not([ng-disabled])')
            .prop('disabled', isDisabled);
  };
  
  var customDisableHandler;
  
  this.onDisable = function (customHandler) {
    customDisableHandler = customHandler;
  };
  
  this.$get = function () {
    return {
      disable: function (rootElement, isDisabled) {
        return (customDisableHandler) ? 
               customDisableHandler(rootElement, isDisabled) : 
               defaultDisableHandler(rootElement, isDisabled);
      }
    }
  };
};

angular.module('rcDisabled', [])
.provider('rcDisabled', rcDisabledProvider)
.directive(rcDisabledDirective);

angular.module('rcDisabledBootstrap', ['rcDisabled'])
.provider('rcDisabled', rcDisabledProvider)
.directive(rcDisabledDirective)
.config(['rcDisabledProvider', function(rcDisabledProvider) {
  rcDisabledProvider.onDisable(function(rootElement, isDisabled) {
    var jqElement = jQuery(rootElement);
      
    jqElement = jqElement
                  .find(':not([rc-disabled])')
                  .filter(function(index) {
                    return jQuery(this).parents().not(jqElement).filter('[rc-disabled]').length === 0;
                  })
                  .filter('input:not([ng-disabled]), button:not([ng-disabled]), .btn, li')
                  .add(jqElement);
            
    // if the Bootstrap "Button" jQuery plug-in is loaded, use it on those
    // that have it configured
    if (jqElement.button) {
      jqElement.find('[data-loading-text]').button((isDisabled) ? 'loading' : 'reset');
    }
            
    jqElement.toggleClass('disabled', isDisabled)
    .filter('input, button')
    .prop('disabled', isDisabled);
  });
}]);

/**
 * @ngdoc module
 * @name rcWizard
 *
 * @description
 * Module to encapsulate the rcWizard directive and rcStep directive.
 */

/**
 * @ngdoc directive
 * @name ng.directive:rcWizard
 *
 * @description
 * Configures the specified element as a wizard.  Uses the jQuery Bootstrap Wizard Plug-in
 *
 * @element ANY
 * @param {name} Specifies the name of the wizard which can be used to look at state 
 * information on the scope.
 */
var rcWizardDirective = {
  'rcWizard': function () {
    return {
      restrict: 'A',
      controller: ['$scope', function ($scope) {
        
        var self;
        var wizardElement;
        var wizardOptions = {};
        var steps = [];
        
        this.currentIndex = 0;
        this.firstIndex = 0;
        this.navigationLength = 0;
        
        this.addStep = function (step) {
          
          steps.push(step);
          
          if (!step.element || !step.submitController) return;
          
          // if a rcSubmitController is being used, automatically add a _hidden_ 
          // submit button so that 
          
          // in order to place an submit button that is still functional it 
          // has to technically be visible, so instead we place it way off 
          // screen
          jQuery(step.element)
            .append('<input type="submit" tabindex="-1" style="position: absolute; left: -9999px; width: 1px; height: 1px;" />')
            .attr('action', 'javascript:void(0);');
          
          // bind to the submit complete event on the rcSubmitController and 
          // if the action was successful, trigger a next on the wizard.
          step.submitController.onSubmitComplete(function (evt) {
            if (evt.success) {
              onForward(step);
            }
          });
        };
          
        this.forward = function () {
          
          if (steps.length)
          
          var currentStep = (steps.length > self.currentIndex) ? steps[self.currentIndex] : null;
          
          if (0 < steps.length && !currentStep) return;
          
          if (0 < steps.length && currentStep.submitController) {
            currentStep.submitController.submit();
          } else {
            onForward(currentStep);
          }
        };
        
        var onForward = function(currentStep) {
          
          if (0 < steps.length && 
            currentStep.formController && 
            currentStep.formController.$invalid) return;
          
          wizardElement.bootstrapWizard('next');
        };
        
        this.backward = function () {
          wizardElement.bootstrapWizard('previous');
        };
        
        var onTabChange = function (activeTab, navigation, currentIndex, nextTab) {
          
          self.currentIndex = nextTab;
          self.firstIndex = wizardElement.bootstrapWizard('firstIndex');
          self.navigationLength = wizardElement.bootstrapWizard('navigationLength');
          
          if (!$scope.$$phase) $scope.$apply();
        };
        
        var onTabClick = function (activeTab, navigation, currentIndex, clickedIndex) {
            return false;
        };
        
        var onTabShow = function (activeTab, navigation, currentIndex) {
          
          if (currentIndex > 0) {
            wizardElement
              .find('.nav li:gt(' + (currentIndex - 1) + ')')
              .removeClass("success");
            wizardElement.find('.nav li:lt(' + currentIndex + ')')
              .addClass("success");
          } else {
            wizardElement.find('.nav li').removeClass("success");
          }
          
          // if a rcStep is being used on the current tab, 
          // automatically focus on the first input of the current tab. This
          // allows for easier keyboard-ony navigation.
          if (steps.length > currentIndex && steps[currentIndex].element) {
            steps[currentIndex].element.find('input').first().focus();
          }
        };
        
        var updateWizard = function (options) {
          
          wizardOptions = options;
          
          if (wizardElement) {
            wizardElement.bootstrapWizard(options);
            self.currentIndex = wizardElement.bootstrapWizard('currentIndex');
            self.firstIndex = wizardElement.bootstrapWizard('firstIndex');
            self.navigationLength = wizardElement.bootstrapWizard('navigationLength');
            
            if (!$scope.$$phase) $scope.$apply();
          }
        };
        
        this.setWizardElement = function (element) {
          
          wizardElement = element;
          self = this;
          updateWizard({
            'onTabChange': onTabChange,
            'onTabShow': onTabShow,
            'onTabClick': onTabClick
          });
        };
      }],
      compile: function (cElement, cAttributes, transclude) {
        return {
          pre: function (scope, formElement, attributes, wizardController) {
            // put a reference to the wizardcontroller on the scope so we can 
            // use some of the properties in the markup
            scope.rc = scope.rc || {};
            scope.rc[attributes.rcWizard] = wizardController;
          },
          post: function (scope, element, attributes, wizardController) {
            // let the controller know about the element
            wizardController.setWizardElement(element);
            if (!scope.$$phase) scope.$apply();
          }
        };
      }
    }
  }
};

/**
 * @ngdoc directive
 * @name ng.directive:rcStep
 *
 * @description
 * Configures the specified element as a wizard-step.  Tells the parent rcWizard
 * controller about the step including any optional controllers.
 *
 * @element ANY
 * @param NONE
 */
var rcWizardStepDirective = {
  'rcStep': function () {
    return {
      restrict: 'A',
      require: ['^rcWizard', '?form', '?rcSubmit'],
      link: function (scope, element, attributes, controllers) {
        
        var wizardController = controllers[0];
        
        // find all the optional controllers for the step
        var formController = controllers.length > 1 ? controllers[1] : null;
        var submitController = controllers.length > 2 ? controllers[2] : null;
        
        // add the step to the wizard controller
        var step = wizardController.addStep({ 
          'element': element, 
          'attributes': attributes, 
          'formController': formController,
          'submitController': submitController });
      }
    };
  }
};

angular.module('rcWizard', ['ng'])

.directive(rcWizardDirective)
.directive(rcWizardStepDirective);

// define controller for wizard
var SampleWizardController = ['$scope', '$q', '$timeout',
function ($scope, $q, $timeout) {
  
  $scope.user = {};
  
  $scope.saveState = function() {
    var deferred = $q.defer();
    
    $timeout(function() {
      deferred.resolve();
    }, 5000);
    
    return deferred.promise;
  };
  
  $scope.completeWizard = function() {
    alert('Completed!');
  }
}];
