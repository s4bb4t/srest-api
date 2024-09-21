document.addEventListener('DOMContentLoaded', () => {
  const taskInput = document.getElementById('taskInput');
  const taskList = document.getElementById('taskList');
  const addTaskButton = document.getElementById('addTask');
  const taskForm = document.getElementById('taskForm');

  // Fetch all tasks from the API
  function fetchTasks() {
      fetch('https://easydev.club/api/v1/todos')
          .then(response => response.json())
          .then(data => {
              data.forEach(task => {
                  addTaskToUI(task);
              });
          });
  }

  // Add task to UI
  function addTaskToUI(task) {
      const li = document.createElement('li');
      li.textContent = task.title;
      
      const deleteButton = document.createElement('button');
      deleteButton.textContent = 'Delete';
      deleteButton.onclick = () => deleteTask(task.id, li);

      li.appendChild(deleteButton);
      taskList.appendChild(li);
  }

  // Add new task
  taskForm.onsubmit = (e) => {
      e.preventDefault();  // Prevent form submission
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
          taskInput.value = ''; // Clear input
      });
  };

  // Delete task
  function deleteTask(taskId, taskElement) {
      fetch(`https://easydev.club/api/v1/todos/${taskId}`, {
          method: 'DELETE'
      })
      .then(() => {
          taskElement.remove();
      });
  }

  // Fetch tasks on load
  fetchTasks();
});
