class Solution {
public:
    bool isPalindrome(ListNode* head) {
        ListNode* slow {head};
        ListNode* fast {head};
        ListNode* next;
        ListNode* prev {new ListNode()};

        while(fast && fast->next) {
            slow = slow->next;
            fast = fast->next->next; 

            next = head->next;
            head->next = prev;
            prev = head;
            head = next;
        }
    
        if (fast) {
            slow = slow->next;
        }
        head = prev;

        while (slow) {
            if (head->val != slow->val) {
                return false;
            }
            head = head->next;
            slow = slow->next;
        }
        return true;
    }
};