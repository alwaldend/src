class Solution {
    public List<String> letterCombinations(String digits) {
        if(digits.length() == 0){
            List<String> result = new ArrayList<>();
            return result;
        }
        List<String> res = combine(digits); 
        return res;
    }
    
    public List<String> combine(String digit){
        if(digit.length() == 0 ){
            List<String> result = new ArrayList<>();
            result.add("");
            return result;
        }
        
        String[] codes = {"","","abc","def","ghi","jkl","mno","pqrs","tuv","wxyz"};
        
        char c = digit.charAt(0);
        
        String  digits_left = digit.substring(1);
        
        List<String> res = combine(digits_left);
        
        List<String> result = new ArrayList<>();
        
        String code_for_current_digit = codes[c-'0'];
        
        for(int i=0;i<code_for_current_digit.length();i++){
            char code_char = code_for_current_digit.charAt(i);
            
            if(!res.isEmpty()){
                for(String s : res){
                    result.add(code_char + s);
                }    
            }
            else{
                res.add(String.valueOf(code_char));
            }
            
        }
        
        
        return result;
    }
}