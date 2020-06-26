.text
    bras clr_test_run
    .asciz "Test verschiedener Befehle"
clr_test_text:
    .asciz "clr test"
    .even

clr_test_run:
    # clr
    lea     clr_test_text(%pc), %a6
    moveq   #1,%d0
    clrw    %d0
    tstw    %d0
    bnes    error
    moveq   #-1, %d0
    clrw    %d0
    swap    %d0
    tstw    %d0
    beqs    error
success:
    trap    #0
error:
    trap    #1
    