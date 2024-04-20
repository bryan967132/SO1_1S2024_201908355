from locust import HttpUser, task, between
import json

class Reader:
    def __init__(self) -> None:
        self.index = 0

    def load(self):
        self.__data = json.loads(open('./data.json', 'r', encoding='utf-8').read())

    def pickNext(self):
        data = self.__data[self.index]
        self.index = self.index + 1 if self.index < 999 else 0
        return data

class MusicUser(HttpUser):
    wait_time = between(1, 5)

    reader = Reader()
    reader.load()

    def on_start(self):
        print('>>> INICIANDO ENVIO DE TRÁFICO')

    @task
    def sendMusicInfo(self):
        print('ENVIA')
        headers = {'content-type': 'application/json'}
        payload = json.dumps(self.reader.pickNext())
        print(payload)
        index = index + 1 if index < 999 else 0
        self.client.post("/sendMusicInfo", json=payload, headers=headers)