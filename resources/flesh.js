var get_test_text = function() {
  test_text = $('#test-text').val();
  return test_text == "" ? null : test_text;
};


var get_grade_level = function(test_text) {
  $.ajax({
    url: '/flesh',
    type: 'POST',
    dataType: 'json',
    contentType: 'application/json; charset=utf-8',
    data: JSON.stringify({'text': test_text}),
    timeout: 3000,
    success: function(json) {
      $('#flesch-kincaid-grade').text(json['grade']);
      $('#flesch-score').text(json['score']);
    }
  });
};


var set_form_actions = function() {
  $("#test-text-form").submit(function(e) {
      e.preventDefault();
  });

  $('#submit').click(function () {
      text_text = get_test_text();
      if (text_text != null) {
        get_grade_level(test_text);
      }
  });
};


$(function () {
    set_form_actions();
});
