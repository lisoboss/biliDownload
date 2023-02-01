package tools

import (
	"bili/tools/media"
	"os"
)

func MediaMerge(videoPath, audioPath, savePath string) (err error) {
	rVideo, err := media.NewRoot(videoPath)
	if err != nil {
		return
	}
	rAudio, err := media.NewRoot(audioPath)
	if err != nil {
		return
	}

	// TrackExtends
	rvTrackExtends := rVideo.GetBox("trex")
	raTrackExtends := rAudio.GetBox("trex")

	// rVideo add rAudio 'trex'
	raTrackExtendsLength := raTrackExtends.Length()
	_pBox := rvTrackExtends.ParentBox()
	_pBox.AddContainerBox(raTrackExtends)
	_pBox.AddLength(raTrackExtendsLength)

	// Track
	rvTrack := rVideo.GetBox("trak")
	raTrack := rAudio.GetBox("trak")

	// rVideo add rAudio 'trak'
	raTrackLength := raTrack.Length()
	_pBox = rvTrack.ParentBox()
	_pBox.AddContainerBox(raTrack)
	_pBox.AddLength(raTrackLength)

	// new media
	newMedia := new(media.Root)

	// add box
	rvIndex := 0
	rvBoxList := rVideo.BoxList()
	rvBoxListLen := len(rvBoxList)
	for i := 0; i < rvBoxListLen; i++ {
		box := rvBoxList[i]
		if box.Type() == "moof" {
			rvIndex = i
			break
		}
		newMedia.AddBox(box)
	}

	// raIndex
	raIndex := 0
	raBoxList := rAudio.BoxList()
	raBoxListLen := len(raBoxList)
	for i := 0; i < raBoxListLen; i++ {
		box := raBoxList[i]
		if box.Type() == "moof" {
			raIndex = i
			break
		}
	}

	// add rAudio audio data(moof+mdat)
	var sequenceNumber uint32 = 1
	for rvIndex < rvBoxListLen || raIndex < raBoxListLen {
		if rvIndex < rvBoxListLen {
			box := rvBoxList[rvIndex]
			_mbox := box.GetBox("mfhd")
			_movieFragmentHeaderBox := new(media.MovieFragmentHeaderBox)
			_movieFragmentHeaderBox.Load(_mbox.Dump())
			_movieFragmentHeaderBox.SetSequenceNumber(sequenceNumber)
			sequenceNumber++
			_mbox.Load(_movieFragmentHeaderBox.Dump())
			newMedia.AddBox(box)
			rvIndex++
			newMedia.AddBox(rvBoxList[rvIndex])
			rvIndex++
		}
		if raIndex < raBoxListLen {
			box := raBoxList[raIndex]
			_mbox := box.GetBox("mfhd")
			_movieFragmentHeaderBox := new(media.MovieFragmentHeaderBox)
			_movieFragmentHeaderBox.Load(_mbox.Dump())
			_movieFragmentHeaderBox.SetSequenceNumber(sequenceNumber)
			sequenceNumber++
			_mbox.Load(_movieFragmentHeaderBox.Dump())
			newMedia.AddBox(box)
			raIndex++
			newMedia.AddBox(raBoxList[raIndex])
			raIndex++
		}
	}

	// save newMedia
	fileOut, err := os.Create(savePath)
	if err != nil {
		return
	}
	defer func() {
		_ = fileOut.Close()
	}()
	_, err = fileOut.Write(newMedia.Dump())
	return
}
