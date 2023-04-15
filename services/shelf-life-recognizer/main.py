import easyocr
import re

from flask import Flask, jsonify, request, json

app = Flask(__name__)
pattern = r"(\d{2}[\.\,]\d{2}[\.\,]\d{2})"


def index():
    try:
        payload = json.loads(request.data)
        text = text_recognition(payload.get("file_path"))
        datesRegex = re.compile(pattern)
        result = datesRegex.findall(text)
        return jsonify({"dates": result})
    except Exception as err:
        return jsonify({"error": err})


def text_recognition(file_path):
    reader = easyocr.Reader(["ru"])
    result = reader.readtext(file_path, detail=0)
    return " ".join(result).lower()


if __name__ == "__main__":
    app.run(14824)
