package main

var aboutText = `NVTool 2.0

Video encoding tool based on rigaya's NVEncC.

Check out the link below to see if your graphics card supports NVENC.

https://developer.nvidia.com/video-encode-and-decode-gpu-support-matrix-new

Options Reference:

Preset

Encode quality preset.

- P1 (= performance)
- P2
- P3
- P4 (= default)
- P5
- P6
- P7 (= quality)


Quality

Set target quality when using VBR mode. (0.0-51.0, 0 = automatic)


Bitrate

Set bitrate in kbps.

To enable constant quality mode, set the value to 0.


AQ

- temporal - Enable adaptive quantization between frames.

- spatial - Enable adaptive quantization in frame.


Strength

Specify the AQ strength. 1 = weak, 15 = strong, 0 = auto


HEVC

Specify the output codec as HEVC, default is h264.


Resize

Set output resolution. When it is different from the input resolution, HW/GPU resizer (Lanczos 4x4) will be activated automatically.

If not specified, it will be same as the input resolution. (no resize)

*Special Values*

- 0 ... Will be same as input.
- One of width or height as negative value
  Will be resized keeping aspect ratio, and a value which could be divided by the negative value will be chosen.

Example: 
input  1280x720
Resize 1024x576 -> normal
Resize 960x0    -> resize to 960x720 (0 will be replaced to 720, same as input)
Resize 1920x-2  -> resize to 1920x1080 (calculated to keep aspect ratio)


Filter Reference:

KNN

Strong noise reduction filter.

Parameters:
- radius=<int> (default=3, 1-5)
  radius of filter.
- strength=<float> (default=0.08, 0.0 - 1.0)
  Strength of the filter.
- lerp=<float> (default=0.2, 0.0 - 1.0)
  The degree of blending of the original pixel to the noise reduction pixel.
- th_lerp=<float> (default=0.8, 0.0 - 1.0)
  Threshold of edge detection.

PMD

Rather weak noise reduction by modified pmd method, aimed to preserve edge while noise reduction.

Parameters:
- apply_count=<int> (default=2, 1- )
  Number of times to apply the filter.
- strength=<float> (default=100, 0-100)
  Strength of the filter.
- threshold=<float> (default=100, 0-255)
  Threshold for edge detection. The smaller the value is, more will be detected as edge, which will be preserved.

Unsharp

unsharp filter, for edge and detail enhancement.

Parameters:
- radius=<int> (default=3, 1-9)
  radius of edge / detail detection.
- weight=<float> (default=0.5, 0-10)
  Strength of edge and detail emphasis. Larger value will result stronger effect.
- threshold=<float> (default=10.0, 0-255)
  Threshold for edge and detail detection.

EdgeLevel 

Edge level adjustment filter, for edge sharpening.

Parameters:
- strength=<float> (default=5.0, -31 - 31)
  Strength of edge sharpening. Larger value will result stronger edge sharpening.
- threshold=<float> (default=20.0, 0 - 255)
  Noise threshold to avoid enhancing noise. Larger value will treat larger luminance change as noise.
- black=<float> (default=0.0, 0-31)
  strength to enhance dark part of edges.
- white=<float> (default=0.0, 0-31)
  strength to enhance bright part of edges.

License information:

This software is for personal use only and not allowed for any commercial purposes.

This software uses the font licensed by be5invis.

-----------------------------------------------------------------------------------------
    iosevka
-----------------------------------------------------------------------------------------

The font is licensed under SIL OFL Version 1.1.

The support code is licensed under Berkeley Software Distribution license.

---
---

Copyright (c) 2015-2020 Belleve Invis (belleve@typeof.net).

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
* Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
* Neither the name of Belleve Invis nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL BELLEVE INVIS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

-----------------------

---

Copyright 2015-2020, Belleve Invis (belleve@typeof.net).

This Font Software is licensed under the SIL Open Font License, Version 1.1.

This license is copied below, and is also available with a FAQ at:
http://scripts.sil.org/OFL

--------------------------


SIL Open Font License v1.1
====================================================


Preamble
----------

The goals of the Open Font License (OFL) are to stimulate worldwide
development of collaborative font projects, to support the font creation
efforts of academic and linguistic communities, and to provide a free and
open framework in which fonts may be shared and improved in partnership
with others.

The OFL allows the licensed fonts to be used, studied, modified and
redistributed freely as long as they are not sold by themselves. The
fonts, including any derivative works, can be bundled, embedded, 
redistributed and/or sold with any software provided that any reserved
names are not used by derivative works. The fonts and derivatives,
however, cannot be released under any other type of license. The
requirement for fonts to remain under this license does not apply
to any document created using the fonts or their derivatives.


Definitions
-------------

"Font Software" refers to the set of files released by the Copyright
Holder(s) under this license and clearly marked as such. This may
include source files, build scripts and documentation.

"Reserved Font Name" refers to any names specified as such after the
copyright statement(s).

"Original Version" refers to the collection of Font Software components as
distributed by the Copyright Holder(s).

"Modified Version" refers to any derivative made by adding to, deleting,
or substituting -- in part or in whole -- any of the components of the
Original Version, by changing formats or by porting the Font Software to a
new environment.

"Author" refers to any designer, engineer, programmer, technical
writer or other person who contributed to the Font Software.


Permission & Conditions
------------------------

Permission is hereby granted, free of charge, to any person obtaining
a copy of the Font Software, to use, study, copy, merge, embed, modify,
redistribute, and sell modified and unmodified copies of the Font
Software, subject to the following conditions:

1. Neither the Font Software nor any of its individual components,
   in Original or Modified Versions, may be sold by itself.

2. Original or Modified Versions of the Font Software may be bundled,
   redistributed and/or sold with any software, provided that each copy
   contains the above copyright notice and this license. These can be
   included either as stand-alone text files, human-readable headers or
   in the appropriate machine-readable metadata fields within text or
   binary files as long as those fields can be easily viewed by the user.

3. No Modified Version of the Font Software may use the Reserved Font
   Name(s) unless explicit written permission is granted by the corresponding
   Copyright Holder. This restriction only applies to the primary font name as
   presented to the users.

4. The name(s) of the Copyright Holder(s) or the Author(s) of the Font
   Software shall not be used to promote, endorse or advertise any
   Modified Version, except to acknowledge the contribution(s) of the
   Copyright Holder(s) and the Author(s) or with their explicit written
   permission.

5. The Font Software, modified or unmodified, in part or in whole,
   must be distributed entirely under this license, and must not be
   distributed under any other license. The requirement for fonts to
   remain under this license does not apply to any document created
   using the Font Software.



Termination
-----------

This license becomes null and void if any of the above conditions are
not met.


    DISCLAIMER
    
    THE FONT SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
    EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO ANY WARRANTIES OF
    MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT
    OF COPYRIGHT, PATENT, TRADEMARK, OR OTHER RIGHT. IN NO EVENT SHALL THE
    COPYRIGHT HOLDER BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
    INCLUDING ANY GENERAL, SPECIAL, INDIRECT, INCIDENTAL, OR CONSEQUENTIAL
    DAMAGES, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
    FROM, OUT OF THE USE OR INABILITY TO USE THE FONT SOFTWARE OR FROM
    OTHER DEALINGS IN THE FONT SOFTWARE.

-----------------------------------------------------------------------------------------
    tamzen-font
-----------------------------------------------------------------------------------------

       _____                                  __             _
      |_   _|_ _ _ __ ___  ___ _   _ _ __    / _| ___  _ __ | |_
        | |/ _ | '_  _ \/ __| | | | '_ \  | |_ / _ \| '_ \| __|
        | | (_| | | | | | \__ \ |_| | | | | |  _| (_) | | | | |_
        |_|\__,_|_| |_| |_|___/\__, |_| |_| |_|  \___/|_| |_|\__|
                               |___/


Copyright 2010 Scott Fial <http://www.fial.com/~scott/tamsyn-font/>

Tamsyn font is free.  You are hereby granted permission to use, copy, modify,
and distribute it as you see fit.

Tamsyn font is provided "as is" without any express or implied warranty.

The author makes no representations about the suitability of this font for
a particular purpose.

In no event will the author be held liable for damages arising from the use
of this font.

       _____                                __             _
      |_   _|_ _ _ __ ___  _______ _ __    / _| ___  _ __ | |_
        | |/ _ | '_  _ \|_  / _ \ '_ \  | |_ / _ \| '_ \| __|
        | | (_| | | | | | |/ /  __/ | | | |  _| (_) | | | | |_
        |_|\__,_|_| |_| |_/___\___|_| |_| |_|  \___/|_| |_|\__|


Copyright 2011 Suraj N. Kurapati <https://github.com/sunaku/tamzen-font>

Tamzen font is free.  You are hereby granted permission to use, copy, modify,
and distribute it as you see fit.

Tamzen font is provided "as is" without any express or implied warranty.

The author makes no representations about the suitability of this font for
a particular purpose.

In no event will the author be held liable for damages arising from the use
of this font.

-----------------------------------------------------------------------------------------
    github.com/AllenDang/giu
-----------------------------------------------------------------------------------------

MIT License

Copyright (c) 2019 Allen Dang

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

-----------------------------------------------------------------------------------------
    github.com/go-gl/glfw
-----------------------------------------------------------------------------------------

Copyright (c) 2012 The glfw3-go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

-----------------------------------------------------------------------------------------
    github.com/sqweek/dialog
-----------------------------------------------------------------------------------------

ISC License

Copyright (c) 2018, the dialog authors.

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

-----------------------------------------------------------------------------------------
    github.com/fsnotify/fsnotify
-----------------------------------------------------------------------------------------

Copyright (c) 2012 The Go Authors. All rights reserved.
Copyright (c) 2012-2019 fsnotify Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

------------------------------------------------------------------------------------------

@Eva1ent 2020
Telegram: https://t.me/Eva1ent
`
