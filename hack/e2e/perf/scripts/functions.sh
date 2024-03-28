# functions.sh

function update_json_with_workers_and_time() {
    # The first parameter is the number of workers
    workers=$1
    # The second parameter is the time
    time=$2
    # The third parameter is the path to the JSON file
    JSON_FILE=$3

    jq --argjson newEntry "{\"workers\": $workers, \"time\": \"$time\"}" '. += [$newEntry]' "$JSON_FILE" > "tmp.$$" && mv "tmp.$$" "$JSON_FILE"
}