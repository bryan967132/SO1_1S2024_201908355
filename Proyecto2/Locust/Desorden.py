import json
import random

with open("./Locust/dataordenada.json", "r", encoding="utf-8") as json_file:
    data = json.load(json_file)
    random.shuffle(data)
    with open('./Locust/data.json', 'w', encoding='utf-8') as json_file:
        json.dump(data, json_file, indent=4, ensure_ascii=False)