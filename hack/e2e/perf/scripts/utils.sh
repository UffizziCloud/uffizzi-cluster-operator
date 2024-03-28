# utils.sh

# Source the script containing your function
source ./hack/e2e/perf/scripts/functions.sh

# Store the function name in a variable and shift the arguments
func_name=$1
shift

# Dynamically call the function based on the first argument
if declare -f "$func_name" > /dev/null
then
    # Call the function with the remaining arguments
    "$func_name" "$@"
else
    echo "Function '$func_name' not found"
fi
