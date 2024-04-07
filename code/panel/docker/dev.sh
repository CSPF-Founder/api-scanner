#!/bin/bash

# Start air in the background
air &

# Change to the directory containing your Node.js project
cd frontend

# Start the npm development server
npm run dev
