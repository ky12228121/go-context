package channel

import (
	"fmt"
	"sync"
)

// sync.WaitGroupはすべてのゴルーチンが終了するまで待機することができる
// Add()に指定した数だけDone()が呼び出されるのをWait()で待つ
var wg sync.WaitGroup

func generator(done chan struct{}, num int) <-chan int {
	// makeは参照型(slice, map, channel)を生成、初期化するために使われる
	// newはプリミティブ型を生成、初期化するために使われる
	out := make(chan int)

	// go文を使うと関数呼び出しをゴルーチンとして実行する
	go func() {
		// defer文は遅延実行される
		// defer文の引数は即時評価されるためdefer文の後に引数の値を書き換えてもdefer呼び出し時は変わっていない
		defer wg.Done()

		// forなどにラベルをつけることができ、ラベルを指定してbreakなどが行える
	LOOP:
		// 条件式を書かないforは無限ループ
		for {
			// select文はchannelの動きを制御する
			// channelのみ指定できるswitch文みたいなもの
			// caseでいずれかのchannelが準備(受信 or Clone)できるまで待機する
			// 評価は上からではなく同時、実行可能なchannelが複数ある場合はランダム選択
			select {
			// doneを受信すると動く
			case <-done:
				break LOOP
			// channelが書き込み可能であれば準備されたとして実行される
			// 書き込み可能かどうかはchannelのバッファが埋まってるかどうかで判断され、埋まっている場合は書き込み可能になるまで待機する
			case out <- num:
			}
		}

		// closeでchannelを閉じることができる
		// closeしたchannelには送信できずpanicになり、受信はゼロ値が即時帰ってくる
		// closeされたかどうかはchannel生成時の第二戻り値で判断できる
		close(out)
		// 標準出力はPrint、Printf、Printlnの3つがある
		// Printは最後に改行がつかない出力
		// Printfは書式付き文字列が使える
		// Printlnは最後に改行が入る
		fmt.Println("generator closed")
	}()
	return out
}

func main() {
	done := make(chan struct{})
	gen := generator(done, 1)

	wg.Add(1)

	for i := 0; i < 5; i++ {
		fmt.Println(<-gen)
	}
	close(done)

	// ゴルーチンの終わりを待つ
	// goはmainが終わると他のゴルーチンが動いていてもプログラムが終了する
	wg.Wait()
}
