/***                        ***
  * This is part of Shenmeci *
 ***                        ***/

var translate = function() {
  // change to the dictionary tab
  $('.navbar .nav li a:first').trigger('click');
  var q = $('#q').val();
  $('#results table').addClass('hide');

  var opts = {
    lines: 13, // The number of lines to draw
    length: 7, // The length of each line
    width: 4, // The line thickness
    radius: 10, // The radius of the inner circle
    rotate: 0, // The rotation offset
    color: '#000', // #rgb or #rrggbb
    speed: 1, // Rounds per second
    trail: 60, // Afterglow percentage
    shadow: false, // Whether to render a shadow
    hwaccel: false, // Whether to use hardware acceleration
    className: 'spinner', // The CSS class to assign to the spinner
    zIndex: 2e9, // The z-index (defaults to 2000000000)
    top: 'auto', // Top position relative to parent in px
    left: 'auto' // Left position relative to parent in px
  };
  var spinner = new Spinner(opts).spin($('#results').get(0));
  var spinner_wrapper = $('#results .spinner').wrap('<div class="well-large" />').parent();

  $.getJSON('segment?q=' + q, function(data) {
    var items = [];
    $.each(data.r, function(idx, item) {
      items.push('<tr><td style="white-space: nowrap;">' + item['z'] +
                 '</td><td><ol><li>' +
                 item['m'].split('/').filter(Boolean).join('</li><li>') +
                 '</li></ol></td></tr>');
    });
    if ($.isEmptyObject(data.r)) {
      items.push('<tr><td colspan="2">I\'m speechless.</td></tr>');
    }
    $('#results table tbody').html(items.join(''));
  })
  .error(function() {
    $('#results table tbody').html('<tr><td colspan="2">Something went wrong :-(</td></tr>');
  })
  .complete(function() {
    spinner.stop();
    spinner_wrapper.remove();
    $('#results table').removeClass('hide');
  });
  return false;
};

var getURLParameter = function (key) {
  var regexp = new RegExp('[?|&]' + key + '=' + '([^&;]+?)(&|#|;|$)');
  var value = (regexp.exec(location.search) || [, ""])[1];
  return decodeURIComponent(value.replace(/\+/g, '%20')) || null;
};

// Automatically fill the form and submit when the query parameter is in the URL.
$(function () {
  var q = getURLParameter("q");
  if (q) {
    $('#q').val(q);
    translate();
  }
});
