    .set TrapBase, 32<<4
    .set Trap0,33<4
    .set Trap1, 34<<4

init_test_framwwork:
    lea trap0success(%pc), %a6
    movel %a6, Trap0
    lea trap1error(%pc), %a6
    movel %a6, Trap1
    bras start_test
trap0success:
    moveq  #0, %d0
    moveal %a6, %a0
    stop #2700
trap1error:
    moveq #1, %d0
    moveal %a6, %a0
    stop #2700
 start_test:
    moveq #0, %d0
    moveq #1, %d1
    moveq #2, %d2
    moveq #3, %d3
    moveq #4, %d4
    moveq #5, %d5
    moveq #6, %d6
    moveq #7, %d7

