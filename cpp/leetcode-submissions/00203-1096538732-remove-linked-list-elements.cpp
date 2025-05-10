/**
 * Definition for singly-linked list.
 * struct ListNode {
 *     int val;
 *     ListNode *next;
 *     ListNode() : val(0), next(nullptr) {}
 *     ListNode(int x) : val(x), next(nullptr) {}
 *     ListNode(int x, ListNode *next) : val(x), next(next) {}
 * };
 */
class Solution {
public:
    ListNode* removeElements(ListNode* head, int val) {
        head = new ListNode(-1, head);
        auto ans = head;
        while (head) {
            auto next = head->next;
            if (next && next->val == val) {
                head->next = next->next;
            } else {
                head = head->next;
            }
        }
        return ans->next;
    }
};