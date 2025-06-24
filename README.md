# Run

```
go run main.go
```

# Router

/localhost/tasks : 

* Description: Submit Tasks
* http method: POST
* request body: ``[5, 10, 15]``
* response:
  * ```
    {
        "job_id": "b486d5ff",
        "message": "Task submitted successfully",
        "task_count": 16
    }
    ```

/localhost/status: 

* Description: Get Taks Status
* http method: GET
* response:
  * ```
    {
        "current_time": 179,
        "schedule_history": [
            {
                "time": 0,
                "task_indexes": [
                    0,
                    1
                ],
                "remaining_times": [
                    0,
                    1
                ]
            },
            {
                "time": 1,
                "task_indexes": [
                    1,
                    2
                ],
                "remaining_times": [
                    0,
                    96
                ]
            }
        ]
    }
    ```

/localhost//scheduler:

* Description: Switch Scheduler
* http method: POST
* request:
  * ```
    {
        "strategy": "SRTF"
    }
    ```
* response:
  * ```
    {
        "current_strategy": "SRTF",
        "message": "Scheduler strategy switched to: SRTF"
    }
    ```
