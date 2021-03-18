(function($) {
    'use strict';
    $(function() {
        var todoListItem = $('.todo-list');
        var todoListInput = $('.todo-list-input');
        $('.todo-list-add-btn').on("click", function(event) {
            event.preventDefault();
    
            var item = $(this).prevAll('.todo-list-input').val();
    
            if (item) {
                $.post("/todos", {name:item}, addItem) // 리스트를 입력하고 버튼을 누르면 todos 서버에 name에 item을 넣어서 post로 보내고 응답이 오면 addItem을 호출해서 화면에 추가
                
                //todoListItem.append("<li><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" + item + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>"); -> 서버에 보내지않고 그냥 바로 add
                todoListInput.val("");
            }
        });
    
        var addItem = function(item) {
            if (item.completed) { // item.id가 있어야 remove요청을 할 수 있으므로 추가
                todoListItem.append("<li class='completed'"+ " id = '" + item.id + "'><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' checked='checked' />" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
            } else {
                todoListItem.append("<li "+ " id = '" + item.id + "'><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
            }
        };
    
        $.get('/todos', function(items) {
            items.forEach(e => {
                addItem(e)
            });
        });
    
        todoListItem.on('change', '.checkbox', function() {
            if ($(this).attr('checked')) {
                $(this).removeAttr('checked');
            } else {
                $(this).attr('checked', 'checked');
            }
    
            $(this).closest("li").toggleClass('completed');
        });
    
        todoListItem.on('click', '.remove', function() {
            $(this).parent().remove();
        });
    
    });
    })(jQuery);