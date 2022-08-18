this was ported from https://github.com/aws-samples/amazon-chime-media-capture-pipeline-demo/blob/main/src/processLambda/app/app.py

quick and dirty non prod ready ;)

put audio chunks in ./data/ directory (download from aws e.g. `aws s3 cp s3://tutoring-mediapipelines-chime-test-dev/captures/7541dd5f-4659-4a7e-92e0-3d5d03297979/audio/ ./data/ --recursive`)
run `go run main.go`

