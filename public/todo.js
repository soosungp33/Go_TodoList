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
            var id = $(this).closest("li").attr('id') // 가장 가까운 li 태그의 id를 가져온다.
            var $self = $(this);
            var complete = true;
            if ($(this).attr('checked')) { // 이미 체크되어있으면 false로 변경해준다.
                complete = false;
            }
            // 이미 체크되어있는 상태면 해제시키라는 의미에서 false를 날려주고, 해제되어 있는 상태면 체크하라는 의미에서 true를 날려준다.
            $.get("complete-todo/"+id+"?complete="+complete, function(data) { // 체크되어있는지 아닌지에 따라 다른 응답이 온다.
                if (complete) {
                    $self.attr('checked', 'checked');
                } else {
                    $self.removeAttr('checked');
                }
                $self.closest("li").toggleClass('completed');

                /* -> 서버에 날리지않고 그냥 체크되고 해제됨
                if ($self.attr('checked')) {
                    $self.removeAttr('checked');
                } else {
                    $self.attr('checked', 'checked');
                }
        
                $(this).closest("li").toggleClass('completed');
                */
            })
        });
    
        todoListItem.on('click', '.remove', function() {
            // 서버로 요청을 날린다음에 요청이 처리된 다음 지워야한다.
            // url: todos/id, method: DELETE
            var id = $(this).closest("li").attr('id') // 가장 가까운 li 태그의 id를 가져온다.
            var $self = $(this);
            $.ajax({
                url: "todos/" + id,
                type: "DELETE",
                success: function(data) {
                    if (data.success) {
                        $self.parent().remove(); // 응답이 성공적으로 오면 지워준다.
                    }
                }
            })
            // $(this).parent().remove(); -> 바로 지워짐
        });
    
    });
    })(jQuery);