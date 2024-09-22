document.addEventListener('DOMContentLoaded', () => {
    const taskInput = document.getElementById('taskInput');
    const taskList = document.getElementById('taskList');
    const taskForm = document.getElementById('taskForm');
  
    // Fetch all tasks from the API
    function fetchTasks() {
        fetch('https://easydev.club/api/v1/todos')
            .then(response => response.json())
            .then(data => {
                taskList.innerHTML = ''; // Очистить список перед добавлением
                data.data.forEach(task => {
                    addTaskToUI(task);
                });
            });
    }
  
    // Add task to UI
    function addTaskToUI(task) {
        const li = document.createElement('li');
        li.innerHTML = `
          <input type="checkbox" ${task.isDone ? 'checked' : ''} onchange="updateTask(${task.id}, this.checked)">
          <span>${task.title}</span>
          <button onclick="deleteTask(${task.id}, this.parentElement)">Удалить</button>
        `;
        taskList.appendChild(li);
    }
  
    // Add new task
    taskForm.onsubmit = (e) => {
        e.preventDefault();
        const taskTitle = taskInput.value;
        if (taskTitle.trim() === '') return;
  
        fetch('https://easydev.club/api/v1/todos', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ title: taskTitle })
        })
        .then(response => response.json())
        .then(task => {
            addTaskToUI(task);
            taskInput.value = '';
        });
    };
  
    // Delete task
    window.deleteTask = (taskId, taskElement) => {
        fetch(`https://easydev.club/api/v1/todos/${taskId}`, {
            method: 'DELETE'
        })
        .then(() => {
            taskElement.remove();
        });
    };
  
    // Update task status
    window.updateTask = (taskId, isDone) => {
        fetch(`https://easydev.club/api/v1/todos/${taskId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ isDone })
        });
    };
  
    // Fetch tasks on load and update every minute
    fetchTasks();
    setInterval(fetchTasks, 60000);
  });
  