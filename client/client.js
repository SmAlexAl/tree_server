$(function () {
	$('#tree_db').jstree()
	$('#tree_cache').jstree()

	$.post('http://localhost:8082/tree/init', function(data){
		if (data.status == "OK") {
			$('#tree_db').jstree(true).settings.core.data = data.db;
			$('#tree_db').jstree(true).refresh();

			$('#tree_cache').jstree(true).settings.core.data = data.collection;
			$('#tree_cache').jstree(true).refresh();
		}
	});
});


$(document).on('click','#add',  function(){
	ob = $('#tree_cache').jstree('get_selected')

	parentId = ob[0]
	val = $('#update_input').val()

		$.ajax({
    	url: 'http://localhost:8082/tree',
    	type: 'POST',
    	data: JSON.stringify({parentId: parentId, value: val}),
    	success: function(data) {
    		if (data.status == "OK") {
			col = data.collection;

			$('#tree_cache').jstree(true).settings.core.data = col;
			$('#tree_cache').jstree(true).refresh();
		}
		}
	});
});

$(document).on('click','#apply', function(){
	$.post('http://localhost:8082/tree/apply', function(data){
		if (data.status == "OK") {
			$('#tree_db').jstree(true).settings.core.data = data.db;
			$('#tree_db').jstree(true).refresh();

			$('#tree_cache').jstree(true).settings.core.data = data.collection;
			$('#tree_cache').jstree(true).refresh();
	}
	});
});

$(document).on('click','#get',  function(){
	ob = $('#tree_db').jstree('get_selected')

	id = ob[0]
	$.get('http://localhost:8082/tree', {id: id}, function(data){
		if (data.status == "OK") {
			col = data.collection;

			$('#tree_cache').jstree(true).settings.core.data = col;
			$('#tree_cache').jstree(true).refresh();
		}

	});
});

$(document).on('click','#update',  function(){
	ob = $('#tree_cache').jstree('get_selected')

	id = ob[0]
	val = $('#update_input').val()
	$.ajax({
    	url: 'http://localhost:8082/tree/update',
    	type: 'PUT',
    	data: JSON.stringify({id: id, value: val}),
    	success: function(data) {
    		if (data.status == "OK") {
				col = data.collection;

				$('#tree_cache').jstree(true).settings.core.data = col;
				$('#tree_cache').jstree(true).refresh();
			}
		}
	});
});

$(document).on('click','#delete',  function(){
	ob = $('#tree_cache').jstree('get_selected')
	$.ajax({
    	url: 'http://localhost:8082/tree/delete',
    	type: 'PUT',
    	data: JSON.stringify({id: ob[0]}),
    	success: function(data) {
    		if (data.status == "OK") {
				col = data.collection;

				$('#tree_cache').jstree(true).settings.core.data = col;
				$('#tree_cache').jstree(true).refresh();
			}
		}
	});
});

$(document).on('click','#reset',  function(){
	$.ajax({
    	url: 'http://localhost:8082/tree/reset',
    	type: 'PUT',
    	success: function(data) {
    		if (data.status == "OK") {
				$('#tree_cache').empty().jstree('destroy');
				$('#tree_cache').jstree()

				col = data.collection;

				$('#tree_db').jstree(true).settings.core.data = col;
				$('#tree_db').jstree(true).refresh();
			}
		}
	});
});



