#include <iostream>
template<typename T>
class stack {
    private:
        class Node {
        public:
            T val;
            Node* next;

            Node() {

            }

            Node(T val) {
                this->val = val;
            }
        };

    public:
    Node* head;
    Node* tail;
    int size;

    stack() {
        Node loc(0);
        head = &loc;
        std::cout << head;
        tail = head;
        size = 0;
    }

    void push(T val) {
        tail->next = new Node(val);
        tail = tail->next;
        size++;
    }

    T pop() {
        T ret = tail->val;
        Node* s = head;
        for (; s->next != tail; s = s->next) {
        }
        s->next = NULL;
        tail = s;

        size--;
        return ret;
    }

    bool isEmpty() {
        return size == 0;
    }

};

int main() {
    stack<int> st;
    int init[] = {1, 3, 4, 9, 10};
    for (int i : init) {
        st.push(i);
    }

    while (!st.isEmpty()) {
        std::cout << st.pop() << " ";
    }

    return 0;
}