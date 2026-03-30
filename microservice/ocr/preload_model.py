from modelscope.pipelines import pipeline
from modelscope.utils.constant import Tasks


def main() -> None:
    pipeline(Tasks.ocr_recognition, model="xiaolv/ocr_small")
    print("modelscope ocr model cached")


if __name__ == "__main__":
    main()
