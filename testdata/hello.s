.text   
        nop
        bras start
        nop
start:  lea text1(%pc), %a5
        moveq #0, %d0
        tas %d0
        move.b %a0@+, %a1@-
        addql #1, %a0@(100)
        moveb #1, %a1@(100, %d1:w)
        dbra %d0, start
        moveml %d0-%d7/%a0-%a4/%a6, -(%sp)
        bnes start
        rts
text1:  .ascii "hello world"
        .even
        nop
        nop     
        rte
