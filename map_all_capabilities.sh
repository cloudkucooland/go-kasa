#!/bin/bash

# Commands to test
commands=(
    "info" "status" "emeter" "dimmer" "wifi" "gettime" "gettimezone"
    "getsensorroutines" "getmanualaction" "getschedulerules" "bulbstate"
    "getcountdown" "ambient" "lightsensor" "diagnose" "mcudiagnose"
    "onboarding" "btncheck" "testmode" "sefinfo" "firmwarelist"
)

# Function to get all unique (Model, Software, IP) tuples
# Using grep to extract IP addresses
echo "Discovering all devices..."
ips=$(kasa discover | grep -oE '[0-9]{1,3}(\.[0-9]{1,3}){3}')

# Temporary file for mapping
echo "Model | Software | IP" > unique_models.txt

processed_keys=()

for ip in $ips; do
    # Get model and software version for this IP
    info=$(kasa -t 1 info "$ip")
    model=$(echo "$info" | grep "Model:" | awk '{print $2}')
    sw_ver=$(echo "$info" | grep "Software:" | awk '{print $2}')
    
    [ -z "$model" ] && model="Unknown"
    [ -z "$sw_ver" ] && sw_ver="Unknown"
    
    key="${model}_${sw_ver}"
    
    # Check if we already have this model+sw combo
    if [[ ! " ${processed_keys[@]} " =~ " ${key} " ]]; then
        echo "$model | $sw_ver | $ip" >> unique_models.txt
        processed_keys+=("$key")
    fi
done

echo "Testing unique (Model, Software) combinations..."
echo "Model | Software | Command | Result"
echo "--- | --- | --- | ---"

# Read unique models and test
tail -n +2 unique_models.txt | while IFS='|' read -r model sw_ver ip; do
    model=$(echo "$model" | xargs)
    sw_ver=$(echo "$sw_ver" | xargs)
    ip=$(echo "$ip" | xargs)
    
    for cmd in "${commands[@]}"; do
        output=$(kasa -t 1 "$cmd" "$ip" 2>&1)
        status=$?
        
        if [[ $status -eq 0 ]]; then
            result="Success"
        elif [[ "$output" == *"member not support"* ]]; then
            result="Not Supported"
        else
            result="Error: $status"
        fi
        
        echo "$model | $sw_ver | $cmd | $result"
        sleep 0.2
    done
done
