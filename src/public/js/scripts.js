(function($) {
  'use strict';
  $(function() {
    $('.sync_btn').on("click", function(event) {
      event.preventDefault();

      var today = new Date();
      
      $.get("/sync",
      {}, 
      function() {
        
      });
      $('.novacard').css('background-color', '#F78181')
    });
  });
})(jQuery);