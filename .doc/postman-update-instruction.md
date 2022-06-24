# Postman Update Instruction
1. If you add a folder, add the following code to 'Pre-request script' and replace 'TYPE_IN_FOLDER_NAME' to the current folder name.
   ```
   pm.variables.set('folder_name','TYPE_IN_FOLDER_NAME');
   ```
2. If you make change with last item in the folder. Make sure the last item has the following code. Also, if the original last item change to other place, remove the following code.
   ```
   pm.variables.set("test_counter", 0)
   ```