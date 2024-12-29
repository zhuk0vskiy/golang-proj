import locust
from locust import HttpUser, task, between, events, env

import json
import os

os.system("ls")
os.system("pwd")
# f = open("config.json")
# config = json.load(f)

# os.environ["LOCUST_USERS"] = config["settings"]["max_users_count"]
# os.environ["LOCUST_SPAWN_RATE"] = config["settings"]["users_spawn_rate"]
# os.environ["LOCUST_RUN_TIME"] = config["settings"]["app_run_time"]

# f.close()


@events.quitting.add_listener
def _(environment, **kw):
    environment.process_exit_code = 0

class WebsiteTestUser(HttpUser):
    wait_time = between(0.5, 1.0)
    config = None

    @task
    def foo(self):
#         global config
#         model_msg = config["endpoints"]["/predict"]
#         msg = json.dumps(model_msg)
        # self.client.post(model_url, msg)
        self.client.get("/api/v1/test/studios/1")
