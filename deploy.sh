#!/bin/bash
GIT_HASH=$(git rev-parse --short HEAD)
img=robherley/plague-doctor:$GIT_HASH
echo Building Frontend
npm run build --prefix ../frontend
echo Moving Frontend to Backend DIR
rm -rf static/
mv ../frontend/dist static
echo Building $img:
docker build -t $img .
echo Pushing $img:
docker push $img

while true; do
    read -p "Do you want to update the deployment? (y/n)" yn
    case $yn in
        [Yy]* ) kubectl set image deployment/plague-doctor plague-doctor=$img; break;;
        [Nn]* ) exit;;
        * ) echo "Please answer y or n.";;
    esac
done
