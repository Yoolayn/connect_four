require("project").command({
    nargs = 1,
    command = "Docker",
    map = {
        up = "!docker compose up -d",
        down = "!docker compose down",
    }
})

require("project").registers({
    j = [[0wye$a<Space>``<Esc>ijson:""<Esc>i<Esc>pguiwf`<Esc>]],
})
