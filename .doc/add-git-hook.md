# Add Git hook to local environment
1. Change directory to repo's root.
   ```
   cd repo_path
   ```
2. Add pre-commit file and change permission.
   ```
   touch .git/hooks/pre-commit
   chmod +x .git/hooks/pre-commit
   ```
3. Use `vim .git/hooks/pre-commit` to edit the pre-commit file. 
   ```
   #!/bin/sh
   
   PASS=true
   
   # Run Newman
   make newman-test
   if [[ $? != 0 ]]; then
       printf "\t\033[31mNewman\033[0m \033[0;30m\033[41mFAILURE!\033[0m\n"
       PASS=false
   else
       printf "\t\033[32mNewman\033[0m \033[0;30m\033[42mpass\033[0m\n"
   fi
   
   if ! $PASS; then
       printf "\033[0;30m\033[41mCOMMIT FAILED\033[0m\n"
       exit 1
   else
       printf "\033[0;30m\033[42mCOMMIT SUCCEEDED\033[0m\n"
   fi
   
   exit 0
   ```